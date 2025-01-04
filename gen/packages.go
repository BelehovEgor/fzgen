package gen

import (
	"fmt"
	"go/types"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/thepudds/fzgen/gen/internal/mod"
	"golang.org/x/tools/go/packages"
)

// findFuncFlag describes bitwise flags for findFunc.
// TODO: this is a temporary fork from fzgo/fuzz.FindFunc.
type findFuncFlag uint

const (
	flagMultiMatch findFuncFlag = 1 << iota
	flagRequireFuzzPrefix
	flagExcludeFuzzPrefix
	flagRequireExported
)

// findFuncsGrouped searches for requested functions matching a package pattern and func pattern,
// returning them grouped by package.
func findFuncsGrouped(pkgPattern, funcPattern, constructorPattern string, flags findFuncFlag) ([]*mod.Package, error) {
	report := func(err error) ([]*mod.Package, error) {
		return nil, fmt.Errorf("finding funcs: %v", err)
	}

	funcRe, err := regexp.Compile(funcPattern)
	if err != nil {
		return report(err)
	}
	constructorRe, err := regexp.Compile(constructorPattern)
	if err != nil {
		return report(err)
	}

	pkgs, err := findFuncs(pkgPattern, funcRe, constructorRe, nil, flags)
	if err != nil {
		return report(err)
	}

	for _, p := range pkgs {
		// put constructors a deterministic order.
		// TODO: for now, we'll prefer simpler constructors as approximated by length (so 'New' before 'NewSomething').
		sort.Slice(p.Constructors, func(i, j int) bool {
			if len(p.Constructors[i].FuncName) < len(p.Constructors[j].FuncName) {
				return true
			}
			if len(p.Constructors[i].FuncName) == len(p.Constructors[j].FuncName) {
				return p.Constructors[i].FuncName < p.Constructors[j].FuncName
			}
			return false
		})
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].PkgPath < pkgs[j].PkgPath
	})

	return pkgs, nil
}

func findFuncs(
	pkgPattern string,
	funcPattern, conPattern *regexp.Regexp,
	env []string,
	flags findFuncFlag,
) ([]*mod.Package, error) {
	report := func(err error) error {
		return fmt.Errorf("error while loading packages for pattern %v: %v", pkgPattern, err)
	}

	cfg := &packages.Config{
		Mode: packages.NeedCompiledGoFiles |
			packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedName |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
	}
	if len(env) > 0 {
		cfg.Env = env
	}
	pkgs, err := packages.Load(cfg, pkgPattern)
	if err != nil {
		return nil, report(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package load error for package pattern %v", pkgPattern)
	}

	return getPackagesContent(pkgs, env, funcPattern, conPattern, flags)
}

func getPackagesContent(
	pkgs []*packages.Package, env []string, funcPattern, conPattern *regexp.Regexp, flags findFuncFlag) ([]*mod.Package, error) {
	pkgsContent := make([]*mod.Package, 0)

	for _, pkg := range pkgs {
		content, err := getPackageContent(pkg, env, funcPattern, conPattern, flags)

		if err != nil {
			return nil, err
		}

		pkgsContent = append(pkgsContent, content)
	}

	for _, pkgA := range pkgsContent {
		for _, _pkgB := range pkgsContent {
			setStructImplementedInterfaces(pkgA.Structs, _pkgB.Interfaces)
		}
	}

	return pkgsContent, nil
}

func getPackageContent(
	pkg *packages.Package,
	env []string,
	funcPattern, conPattern *regexp.Regexp,
	flags findFuncFlag,
) (*mod.Package, error) {
	pkgDir := ""
	var err error

	targets := make([]*mod.Func, 0)
	constructors := make([]*mod.Func, 0)
	funcs := make([]*mod.Func, 0)
	structs := make([]*mod.Struct, 0)
	interfaces := make([]*mod.Interface, 0)

	for id, obj := range pkg.TypesInfo.Defs {
		if pkgDir == "" {
			pkgDir, err = goListDir(pkg.PkgPath, env)
			if err != nil {
				return nil, err
			}
		}

		switch objType := obj.(type) {
		case *types.TypeName:
			if !obj.Exported() {
				continue
			}

			objNamed := obj.Type().(*types.Named)
			if types.IsInterface(objNamed) {
				interfaces = append(interfaces, &mod.Interface{
					InterfaceName:  id.Name,
					PkgName:        pkg.Name,
					PkgPath:        pkg.PkgPath,
					PkgDir:         pkgDir,
					TypesInterface: objType.Type().Underlying().(*types.Interface),
					TypesNamed:     objNamed,
				})
			} else if structType, ok := objType.Type().Underlying().(*types.Struct); ok {
				structs = append(structs, &mod.Struct{
					StructName:  id.Name,
					PkgName:     pkg.Name,
					PkgPath:     pkg.PkgPath,
					PkgDir:      pkgDir,
					TypesStruct: structType,
					TypesNamed:  objNamed,
				})
			}
		case *types.Func:
			f := mod.Func{
				FuncName:  id.Name,
				PkgName:   pkg.Name,
				PkgPath:   pkg.PkgPath,
				PkgDir:    pkgDir,
				TypesFunc: objType,
			}

			funcs = append(funcs, &f)
			addTarget(&targets, &f, funcPattern, flags)
			addConstructor(&constructors, &f, conPattern)
		}
	}

	return &mod.Package{
		PkgName:      pkg.Name,
		PkgPath:      pkg.PkgPath,
		Targets:      targets,
		Constructors: constructors,
		Funcs:        funcs,
		Structs:      structs,
		Interfaces:   interfaces,
	}, nil
}

func addTarget(targets *[]*mod.Func, f *mod.Func, funcPattern *regexp.Regexp, flags findFuncFlag) error {
	if isInterfaceRecv(f.TypesFunc) {
		// TODO: control via flag?
		// TODO: merge back to fzgo/fuzz.FindFunc?
		return nil
	}
	if flags&flagExcludeFuzzPrefix != 0 && strings.HasPrefix(f.FuncName, "Fuzz") {
		// skip any function that already starts with Fuzz
		return nil
	}
	if flags&flagRequireFuzzPrefix != 0 && !strings.HasPrefix(f.FuncName, "Fuzz") {
		// skip any function that does not start with Fuzz
		return nil
	}
	if flags&flagRequireExported != 0 {
		if !isExportedFunc(f.TypesFunc) {
			return nil
		}
	}

	if funcPattern.MatchString(f.FuncName) {
		// found a match.
		// check if we already found a match in a prior iteration our of chains.
		if len(*targets) > 0 && flags&flagMultiMatch == 0 {
			return fmt.Errorf("multiple matches not allowed. multiple matches for func %v: %v.%v and %v.%v",
				funcPattern, f.PkgPath, f.FuncName, (*targets)[0].PkgPath, (*targets)[0].FuncName)
		}

		*targets = append(*targets, f)
	}

	return nil
}

func addConstructor(constructors *[]*mod.Func, f *mod.Func, constructorRe *regexp.Regexp) {
	if isConstructor(f.TypesFunc) && constructorRe.MatchString(f.FuncName) {
		*constructors = append(*constructors, f)
	}
}

func setStructImplementedInterfaces(structs []*mod.Struct, interfaces []*mod.Interface) {
	for _, _struct := range structs {
		for _, _interface := range interfaces {
			if types.AssignableTo(_struct.TypesNamed, _interface.TypesInterface) {
				_interface.Implementations = append(_interface.Implementations, _struct)
			} else if types.AssignableTo(types.NewPointer(_struct.TypesNamed), _interface.TypesInterface) {
				_interface.Implementations = append(_interface.Implementations, &mod.Struct{
					StructName:  _struct.StructName,
					PkgName:     _struct.PkgName,
					PkgPath:     _struct.PkgPath,
					PkgDir:      _struct.PkgDir,
					TypesStruct: _struct.TypesStruct,
					TypesNamed:  _struct.TypesNamed,
					AsPointer:   true,
				})
			}
		}
	}
}

// goListDir returns the dir for a package import path.
func goListDir(pkgPath string, env []string) (string, error) {
	if len(env) == 0 {
		env = os.Environ()
	}

	// TODO: use build tags, or not?
	// cmd := exec.Command("go", "list", "-f", "{{.Dir}}", buildTagsArg, pkgPath)
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgPath)
	cmd.Env = env
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		// If this fails with a pkgPath as empty string, check packages.Config.Mode
		fmt.Fprintf(os.Stderr, "fzgen: 'go list -f {{.Dir}} %v' failed for pkgPath %q\n%v\n", pkgPath, pkgPath, string(out))
		return "", fmt.Errorf("failed to find directory for package %q: %v", pkgPath, err)
	}
	result := strings.TrimSpace(string(out))
	if strings.Contains(result, "\n") {
		return "", fmt.Errorf("multiple directory results for package %v", pkgPath)
	}
	return result, nil
}

// goList returns a list of packages for a package pattern.
// There is probably a more refined way to do this, but it might do 'go list' anyway.
func goList(dir string, args ...string) ([]string, error) {
	if dir == "" {
		dir = "."
	}

	cmdArgs := append([]string{"list"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		// If this fails with a pkgPath as empty string, check packages.Config.Mode
		fmt.Fprintf(os.Stderr, "fzgen: 'go list' failed for args %q:\n%s\n", args, string(out))
		return nil, fmt.Errorf("failed to find package for args %q: %v", args, err)
	}

	var result []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, line)
	}
	return result, nil
}

// isInModule reports if dir appears to be within a module with a 'go.mod'.
func isInModule(dir string) (bool, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	s := strings.TrimSpace(string(out))
	if s == "" || s == os.DevNull {
		// Go 1.11 reports empty string for no 'go.mod'.
		// Go 1.12+ report os.DevNull for no 'go.mod'
		return false, nil
	}
	return true, nil
}

// TODO: would be good to find some canonical documentation or example of this.
func isExportedFunc(f *types.Func) bool {
	if !f.Exported() {
		return false
	}
	// the function itself is exported, but it might be a method on an unexported type.
	sig, ok := f.Type().(*types.Signature)
	if !ok {
		return false
	}
	recv := sig.Recv()
	if recv == nil {
		// not a method, and the func itself is exported.
		return true
	}

	n, err := namedType(recv)
	if err != nil {
		// don't treat as fatal error.
		fmt.Fprintf(os.Stderr, "fzgen: warning: failed to determine if exported for receiver %v for func %v: %v\n",
			recv, f, err)
		return false
	}

	return n.Obj().Exported()
}

// isInterfaceRecv helps filter out interface receivers such as 'func (interface).Is(error) bool'
// Previously would have issues from errors.Is:
//
//	x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target)
func isInterfaceRecv(f *types.Func) bool {
	sig, ok := f.Type().(*types.Signature)
	if !ok {
		return false
	}
	recv := sig.Recv()
	if recv == nil {
		// not a method
		return false
	}
	// TODO: this might be redundant check to do both Type() and Type().Underlying(), but shouldn't hurt.
	_, ok = recv.Type().(*types.Interface)
	if ok {
		return true
	}
	_, ok = recv.Type().Underlying().(*types.Interface)
	return ok
}

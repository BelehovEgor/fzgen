# FzGen

Автоматически генерирует фаззинг цели для переданного пакета.

## Установка

```
$ go install github.com/BelehovEgor/fzgen/cmd/fzgen@latest
$ go get github.com/BelehovEgor/fzgen/fuzzer
```

## Настройки

[-ctor=<target-constructor-regexp>] - Регулярное выражение, определяющее имя конструкторов, которые будут использоваться для создания типов данных пакета.

[-unexported] - В тестовом файле будут использоваться несэкспортированные структуры, функции.

[-mocks] [-mocksPackagePrefix] [-mocksDepth] - Набор флагов определяющий интеграцию с Mockery (который должен быть установлен для использования). 
- [-mocks] Включение моков в генерации.
- [-mocksPackagePrefix=] Имя пакета, в котором лежит тестовый файл. Необходимо для правильной генерации импортов сгенерированных моков.
- [-mocksDepth=] Максимальная глубина при генерации моков. (Если для генерации мока, требуется сгенерировать еще один мок - глубина увеличивается на 1) Default - 5.

[-structFillMode=] - Выбирает настройку заполнения структур
- "Constructors" - если есть конструктор в приоритете выбирается конструктор для создания структуры
- "Random" - заполнение структур по полям
- "ConstructorsAndRandom" - случайно выбирает как заполнять структуру

[-fillUnexported] - Включает возможность заполнять несэкспортированные поля в структурах.

[-llm=<name-of-client>] - Включает интеграцию с LLM. Сейчас реализованы интеграции с openrouter, groq, gigachat
- "openrouter"
- "gigachat"
- "groq-qwen"
- "groq-deepseek"

Для использования данных интеграций необходимо указать креды через переменные окружения:
- openrouter_api_key
- gigachat_api_key
- groq_api_key

Для задержки между запросами - переменная timeout - в секундах указывает паузу между запросами в клиент LLM.

[-chain] [-parallel] - Наследие от первой версии FzGen, не поддержаны все изменения, нет полного тестирования.

## Использование

```
fzgen {flags} {name_of_package}
```

### Примеры

```
fzgen io 
```

```go
package main

// Edit if desired. Code generated by "fzgen io".

import (
	"fmt"
	io "io"
	"reflect"
	"testing"

	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_N1_LimitedReader_Read(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var l *io.LimitedReader
		var p []byte
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&l, &p)
		if err != nil || l == nil {
			return
		}

		// Put here your precondition of func arguments...

		l.Read(p)

		// Put here your postcondition of func results...
	})
}
//...

func fabric_constructor_wrapper_TeeReader_0(
	r io.Reader,
	w io.Writer,
) (result io.Reader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	res := io.TeeReader(
		r,
		w,
	)
	return res, err
}

func fabric_constructor_wrapper_NopCloser_0(
	r io.Reader,
) (result io.ReadCloser, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	res := io.NopCloser(
		r,
	)
	return res, err
}

func fabric_constructor_wrapper_NewSectionReader_0(
	r io.ReaderAt,
	off int64,
	n int64,
) (result io.SectionReader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	res := io.NewSectionReader(
		r,
		off,
		n,
	)
	return *res, err
}

func fabric_constructor_wrapper_MultiWriter_0(
	writers []io.Writer,
) (result io.Writer, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	res := io.MultiWriter(
		writers...,
	)
	return res, err
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	FabricFuncsForCustomTypes["io.Writer"] = append(FabricFuncsForCustomTypes["io.Writer"], reflect.ValueOf(fabric_constructor_wrapper_MultiWriter_0))
	FabricFuncsForCustomTypes["io.Reader"] = append(FabricFuncsForCustomTypes["io.Reader"], reflect.ValueOf(fabric_constructor_wrapper_TeeReader_0))
	FabricFuncsForCustomTypes["io.WriterAt"] = append(FabricFuncsForCustomTypes["io.WriterAt"], reflect.ValueOf(fabric_interface_io_WriterAt_OffsetWriter))
	FabricFuncsForCustomTypes["io.ReaderAt"] = append(FabricFuncsForCustomTypes["io.ReaderAt"], reflect.ValueOf(fabric_interface_io_ReaderAt_SectionReader))
	FabricFuncsForCustomTypes["io.ReadCloser"] = append(FabricFuncsForCustomTypes["io.ReadCloser"], reflect.ValueOf(fabric_constructor_wrapper_NopCloser_0))
	FabricFuncsForCustomTypes["io.SectionReader"] = append(FabricFuncsForCustomTypes["io.SectionReader"], reflect.ValueOf(fabric_constructor_wrapper_NewSectionReader_0))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_PipeReader))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_OffsetWriter))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_LimitedReader))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_string))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_SectionReader))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_PipeWriter))
	m.Run()
}
```

```
fzgen -unexported io 
```

```go
func Fuzz_N28_pipe_writeCloseError(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var p *io.pipe
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&p)
		if err != nil || p == nil {
			return
		}

		// Put here your precondition of func arguments...

		p.writeCloseError()

		// Put here your postcondition of func results...
	})
}
```

```
fzgen -mocks -mocksPackagePrefix=example -mocksDepth=3 net/http
```

```go
func fabric_mock_interface_21_File(
	t *testing.T,
	error_ error,
	n int,
	err error,
	fileinfo_ []fs.FileInfo,
	error__1 error,
	int64_ int64,
	error__2 error,
	fileinfo__1 fs.FileInfo,
	error__3 error,
) http.File {
	genMock := mocks_2.NewMockFile(t)
	genMock.
		On("Close").
		Return(func() error {
			return error_
		}).
		Maybe()
	genMock.
		On("Read", mock.AnythingOfType("[]byte")).
		Return(func(p []byte) (n int, err error) {
			return n, err
		}).
		Maybe()
	genMock.
		On("Readdir", mock.AnythingOfType("int")).
		Return(func(count int) ([]fs.FileInfo, error) {
			return fileinfo_, error__1
		}).
		Maybe()
	genMock.
		On("Seek", mock.AnythingOfType("int64"), mock.AnythingOfType("int")).
		Return(func(offset int64, whence int) (int64, error) {
			return int64_, error__2
		}).
		Maybe()
	genMock.
		On("Stat").
		Return(func() (fs.FileInfo, error) {
			return fileinfo__1, error__3
		}).
		Maybe()
	return genMock
}
```

```
fzgen -structFillMode=ConstructorsAndRandom net/http
```

```go
func Fuzz_N77_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.ConstructorsAndRandom)
		err := fz.Fill2(&url_0)
		if err != nil {
			return
		}

		// Put here your precondition of func arguments...

		http.Get(url_0)

		// Put here your postcondition of func results...
	})
}
```

```
fzgen -fillUnexported regexp
```

```go
func Fuzz_N39_Regexp_UnmarshalText(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var re *regexp.Regexp
		var text []byte
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors, fuzzer.FillUnexported)
		err := fz.Fill2(&re, &text)
		if err != nil || re == nil {
			return
		}

		// Put here your precondition of func arguments...

		re.UnmarshalText(text)

		// Put here your postcondition of func results...
	})
}
```

```
export openrouter_api_key=key
export timeout=10
fzgen -llm=openrouter sort
```

```go
func Fuzz_N24_SearchFloat64s(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var a []float64
		var x float64
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&a, &x)
		if err != nil {
			return
		}

		// Check if the slice is sorted in ascending order
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return
			}
		}

		sort.SearchFloat64s(a, x)

		// Check if the result index is valid
		result := sort.SearchFloat64s(a, x)
		if result < 0 || result > len(a) {
			t.Error("result index out of bounds")
		}

		// Check if the result index is correct
		if result < len(a) && a[result] < x {
			t.Error("result index is incorrect")
		}
		if result > 0 && a[result-1] >= x {
			t.Error("result index is incorrect")
		}
	})
}
```
# gowcode

Вычислитель выражений для Go. Разбирает и вычисляет выражения с переменными, арифметикой, операторами сравнения, тернарным оператором, литералами списков и словарей, индексацией и богатым набором встроенных функций. Поддерживает пользовательские функции через реестр.

## Содержание

- [Как это работает](#как-это-работает)
- [Быстрый старт](#быстрый-старт)
  - [Eval — одиночное вычисление](#eval--одиночное-вычисление)
  - [Evaluator — переиспользуемый экземпляр](#evaluator--переиспользуемый-экземпляр)
  - [Пользовательские функции](#пользовательские-функции)
- [Переменные](#переменные)
  - [Синтаксис {varname}](#синтаксис-varname)
  - [Управление переменными из выражений](#управление-переменными-из-выражений)
- [Типы](#типы)
- [Операторы](#операторы)
- [Встроенные функции](#встроенные-функции)
  - [Строки](#строки)
  - [Числа](#числа)
  - [Списки](#списки)
  - [Словари](#словари)
  - [Логика](#логика)
  - [Преобразование и проверка типов](#преобразование-и-проверка-типов)
  - [Прочее](#прочее)

---

## Как это работает

```
строка выражения
      │
      ▼
   Lexer  →  токены  ({x} → TokVar)
      │
      ▼
   Parser →  AST
      │
      ▼
  Evaluator
  ├── vars map[string]*value.Value
  └── Registry (встроенные + пользовательские функции)
      │
      ▼
   *value.Value
```

1. **Lexer** (`lexer/`) разбивает строку на токены. `{varname}` распознаётся как отдельный токен переменной.
2. **Parser** (`parser/`) строит AST из потока токенов.
3. **Evaluator** (`eval/`) обходит AST, подставляет переменные из переданного словаря и вызывает функции через **Registry** (`functions/`).
4. Все значения представлены типом `*value.Value` — тегированным объединением, которое хранит `null`, `bool`, `int`, `float`, `string`, `list` или `map`.

---

## Быстрый старт

### Eval — одиночное вычисление

`eval.Eval` — самый простой способ использования. Разбирает и вычисляет выражение за один вызов.

```go
package main

import (
    "fmt"
    "gowcode/eval"
    "gowcode/value"
)

func main() {
    vars := map[string]*value.Value{
        "x": value.IntVal(10),
        "y": value.IntVal(3),
    }

    result, err := eval.Eval("{x} * {y} + 1", vars)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // 31

    // Конкатенация строк
    vars2 := map[string]*value.Value{
        "name": value.StringVal("World"),
    }
    r2, _ := eval.Eval(`"Hello, " + {name} + "!"`, vars2)
    fmt.Println(r2) // Hello, World!

    // Тернарный оператор
    vars3 := map[string]*value.Value{
        "score": value.IntVal(85),
    }
    r3, _ := eval.Eval(`{score} >= 90 ? "A" : {score} >= 75 ? "B" : "C"`, vars3)
    fmt.Println(r3) // B

    // Встроенные функции
    r4, _ := eval.Eval(`upper(trim("  hello  "))`, nil)
    fmt.Println(r4) // HELLO

    // Литерал списка + индексация
    r5, _ := eval.Eval(`[10, 20, 30][1]`, nil)
    fmt.Println(r5) // 20

    // Литерал словаря + индексация
    r6, _ := eval.Eval(`{"name": "Alice", "age": 30}["name"]`, nil)
    fmt.Println(r6) // Alice
}
```

### Evaluator — переиспользуемый экземпляр

Используйте `eval.NewEvaluator`, если нужно вычислять несколько выражений с одним и тем же набором переменных.

```go
package main

import (
    "fmt"
    "gowcode/eval"
    "gowcode/value"
)

func main() {
    vars := map[string]*value.Value{
        "price":    value.FloatVal(9.99),
        "quantity": value.IntVal(5),
        "discount": value.FloatVal(0.1),
    }

    e := eval.NewEvaluator(vars)

    total, _ := e.Evaluate("{price} * {quantity}")
    final, _ := e.Evaluate("{price} * {quantity} * (1 - {discount})")
    label, _ := e.Evaluate(`{final} > 40 ? "expensive" : "cheap"`)

    fmt.Println(total) // 49.95
    fmt.Println(final) // 44.955
    fmt.Println(label) // expensive
}
```

### Пользовательские функции

Создайте `Registry`, зарегистрируйте свои функции и передайте реестр в `NewEvaluatorWithRegistry`.

```go
package main

import (
    "fmt"
    "strings"

    "gowcode/eval"
    "gowcode/functions"
    "gowcode/value"
)

func main() {
    reg := functions.NewRegistry() // включает все встроенные функции

    // Пользовательская функция: greet(name)
    reg.RegisterFunc("greet", func(args []*value.Value) (*value.Value, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("greet: ожидается 1 аргумент")
        }
        return value.StringVal("Привет, " + args[0].AsString() + "!"), nil
    })

    // Пользовательская функция: titleCase(s)
    reg.RegisterFunc("titleCase", func(args []*value.Value) (*value.Value, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("titleCase: ожидается 1 аргумент")
        }
        return value.StringVal(strings.Title(args[0].AsString())), nil
    })

    vars := map[string]*value.Value{
        "user": value.StringVal("alice"),
    }

    e := eval.NewEvaluatorWithRegistry(vars, reg)

    r1, _ := e.Evaluate(`greet(titleCase({user}))`)
    fmt.Println(r1) // Привет, Alice!

    // Встроенные функции можно переопределить
    reg.RegisterFunc("upper", func(args []*value.Value) (*value.Value, error) {
        return value.StringVal("<<" + strings.ToUpper(args[0].AsString()) + ">>"), nil
    })

    r2, _ := e.Evaluate(`upper("world")`)
    fmt.Println(r2) // <<WORLD>>
}
```

> **Примечание:** `NewEvaluator` всегда создаёт реестр со всеми встроенными функциями. `NewEvaluatorWithRegistry` позволяет передать заранее настроенный реестр с добавленными или переопределёнными функциями.

---

## Переменные

### Синтаксис {varname}

Переменные в выражениях обязательно оборачиваются в фигурные скобки: `{имя}`. Это позволяет избежать конфликтов с именами встроенных функций и сделать выражения явными.

```
{price} * {quantity}        // обращение к переменным
upper({name})               // переменная как аргумент функции
{items}[0]                  // индексация по переменной-списку
{config}["timeout"]         // доступ к ключу переменной-словаря
{a} > {b} ? {a} : {b}      // переменные в тернарном операторе
```

Пробелы внутри скобок допускаются: `{ name }` эквивалентно `{name}`.

Фигурные скобки без единственного идентификатора внутри трактуются как литерал словаря:

```
{"key": "value"}            // словарь — не переменная
{a + b}                     // ошибка синтаксиса (не идентификатор)
```

### Управление переменными из выражений

`Evaluator` предоставляет четыре встроенные функции для работы с переменными прямо из выражений. Это удобно для lowcode-сценариев, когда нужно накапливать состояние между несколькими вызовами `Evaluate`.

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `setVar` | `setVar(name, value)` | Записывает переменную в vars; возвращает `value` |
| `getVar` | `getVar(name)` | Читает переменную из vars; возвращает `null`, если нет |
| `deleteVar` | `deleteVar(name)` | Удаляет переменную из vars; возвращает `null` |
| `hasVar` | `hasVar(name)` | Проверяет наличие переменной; возвращает `bool` |

После `setVar` переменная доступна через `{name}` в любом последующем выражении того же экземпляра `Evaluator`.

```go
e := eval.NewEvaluator(nil)

// Записываем переменную из выражения
e.Evaluate(`setVar("total", 0)`)
e.Evaluate(`setVar("total", {total} + 10)`)
e.Evaluate(`setVar("total", {total} + 5)`)

result, _ := e.Evaluate(`{total}`)
fmt.Println(result) // 15

// Проверка наличия
has, _ := e.Evaluate(`hasVar("total")`)
fmt.Println(has) // true

// Удаление
e.Evaluate(`deleteVar("total")`)
has2, _ := e.Evaluate(`hasVar("total")`)
fmt.Println(has2) // false
```

Пример накопления результата в lowcode-цепочке:

```go
e := eval.NewEvaluator(map[string]*value.Value{
    "prices": value.ListVal(
        value.FloatVal(9.99),
        value.FloatVal(4.50),
        value.FloatVal(14.00),
    ),
    "tax": value.FloatVal(0.2),
})

e.Evaluate(`setVar("subtotal", sum({prices}))`)
e.Evaluate(`setVar("taxAmount", {subtotal} * {tax})`)
e.Evaluate(`setVar("total", {subtotal} + {taxAmount})`)

total, _ := e.Evaluate(`{total}`)
fmt.Println(total) // 34.188
```

---

## Типы

| Тип      | Пример литерала            | Go-конструктор            |
|----------|----------------------------|---------------------------|
| `null`   | `null`                     | `value.Nil()`             |
| `bool`   | `true`, `false`            | `value.BoolVal(b)`        |
| `int`    | `42`, `-7`                 | `value.IntVal(i)`         |
| `float`  | `3.14`, `-0.5`             | `value.FloatVal(f)`       |
| `string` | `"hello"`, `'world'`       | `value.StringVal(s)`      |
| `list`   | `[1, 2, 3]`                | `value.ListVal(items...)` |
| `map`    | `{"key": "value"}`         | `value.MapVal()` + `.Set` |

При арифметике и сравнениях выполняется автоматическое приведение типов (например, `int` и `float` → `float`). Оператор `+` для строк выполняет конкатенацию.

---

## Операторы

### Арифметические

| Оператор | Описание                    | Пример              |
|----------|-----------------------------|---------------------|
| `+`      | Сложение / конкатенация     | `2 + 3` → `5`      |
| `-`      | Вычитание                   | `10 - 4` → `6`     |
| `*`      | Умножение                   | `3 * 4` → `12`     |
| `/`      | Деление                     | `7 / 2` → `3`      |
| `%`      | Остаток от деления          | `7 % 3` → `1`      |
| `**`     | Возведение в степень        | `2 ** 8` → `256.0` |

Целочисленное деление усекает дробную часть. `**` всегда возвращает `float`.

### Сравнения

| Оператор | Описание              |
|----------|-----------------------|
| `==`     | Равно                 |
| `!=`     | Не равно              |
| `<`      | Меньше                |
| `<=`     | Меньше или равно      |
| `>`      | Больше                |
| `>=`     | Больше или равно      |

### Логические

| Оператор | Описание         |
|----------|------------------|
| `&&`     | Логическое И     |
| `\|\|`   | Логическое ИЛИ   |
| `!`      | Логическое НЕ    |

### Прочие

| Синтаксис           | Описание                                          |
|---------------------|---------------------------------------------------|
| `{varname}`         | Обращение к переменной                            |
| `cond ? a : b`      | Тернарный оператор                                |
| `list[i]`           | Индексация списка (с 0; отрицательный — с конца)  |
| `map["key"]`        | Доступ к ключу словаря                            |

---

## Встроенные функции

### Строки

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `upper` | `upper(s)` | Перевод в верхний регистр |
| `lower` | `lower(s)` | Перевод в нижний регистр |
| `trim` | `trim(s)` / `trim(s, chars)` | Удалить пробелы или указанные символы с обоих концов |
| `trimLeft` | `trimLeft(s)` / `trimLeft(s, chars)` | Удалить с левого конца |
| `trimRight` | `trimRight(s)` / `trimRight(s, chars)` | Удалить с правого конца |
| `trimPrefix` | `trimPrefix(s, prefix)` | Удалить префикс, если присутствует |
| `trimSuffix` | `trimSuffix(s, suffix)` | Удалить суффикс, если присутствует |
| `replace` | `replace(s, old, new)` / `replace(s, old, new, n)` | Заменить до `n` вхождений (`-1` = все) |
| `replaceAll` | `replaceAll(s, old, new)` | Заменить все вхождения |
| `startsWith` | `startsWith(s, prefix)` | Начинается ли строка с `prefix`; возвращает `bool` |
| `endsWith` | `endsWith(s, suffix)` | Заканчивается ли строка на `suffix`; возвращает `bool` |
| `contains` | `contains(s, sub)` | Содержит ли строка подстроку; возвращает `bool`; работает и для списков |
| `indexOf` | `indexOf(s, sub)` | Первый индекс подстроки / элемента, или `-1` |
| `lastIndexOf` | `lastIndexOf(s, sub)` | Последний индекс подстроки / элемента, или `-1` |
| `split` | `split(s, sep)` | Разбить строку в список |
| `join` | `join(list, sep)` | Объединить элементы списка в строку |
| `len` | `len(s)` | Длина в символах Unicode; работает и для списков и словарей |
| `substr` | `substr(s, start)` / `substr(s, start, end)` | Подстрока с поддержкой Unicode; отрицательные индексы считаются с конца |
| `repeat` | `repeat(s, n)` | Повторить строку `n` раз |
| `padLeft` | `padLeft(s, width)` / `padLeft(s, width, pad)` | Дополнить строку слева до `width` символов |
| `padRight` | `padRight(s, width)` / `padRight(s, width, pad)` | Дополнить строку справа до `width` символов |
| `format` | `format(template, args...)` | Форматирование в стиле Printf (`%s`, `%d`, `%f`, …) |
| `charAt` | `charAt(s, i)` | Символ по индексу `i`; отрицательный — с конца |
| `reverse` | `reverse(s)` | Перевернуть строку или список |

**Примеры:**

```
upper("hello")                         // "HELLO"
trim("  hi  ")                         // "hi"
trim("--hello--", "-")                 // "hello"
replace("aabbcc", "b", "X", 1)         // "aaXbcc"
split("a,b,c", ",")                    // ["a", "b", "c"]
join(["x", "y", "z"], "-")             // "x-y-z"
substr("Hello, World", 7, 12)          // "World"
padLeft("5", 3, "0")                   // "005"
format("Hi %s, you are %d", "Bob", 30) // "Hi Bob, you are 30"
```

---

### Числа

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `abs` | `abs(n)` | Абсолютное значение |
| `floor` | `floor(n)` | Округление вниз |
| `ceil` | `ceil(n)` | Округление вверх |
| `round` | `round(n)` / `round(n, decimals)` | Округление до целого или до `decimals` знаков после запятой |
| `trunc` | `trunc(n)` | Усечение дробной части к нулю |
| `sqrt` | `sqrt(n)` | Квадратный корень |
| `pow` | `pow(base, exp)` | Возведение в степень (результат `float`) |
| `sign` | `sign(n)` | Знак числа: `-1`, `0` или `1` |
| `log` | `log(n)` / `log(n, base)` | Натуральный логарифм или логарифм по основанию `base` |
| `log2` | `log2(n)` | Логарифм по основанию 2 |
| `log10` | `log10(n)` | Логарифм по основанию 10 |
| `pi` | `pi()` | π ≈ 3.14159… |
| `e` | `e()` | Число Эйлера ≈ 2.71828… |
| `isNaN` | `isNaN(n)` | Является ли значение NaN? |
| `isInf` | `isInf(n)` | Является ли значение ±Infinity? |
| `min` | `min(a, b, …)` / `min(list)` | Минимальное значение |
| `max` | `max(a, b, …)` / `max(list)` | Максимальное значение |
| `clamp` | `clamp(n, lo, hi)` | Ограничить `n` диапазоном `[lo, hi]` |
| `sum` | `sum(a, b, …)` / `sum(list)` | Сумма значений |
| `avg` | `avg(a, b, …)` / `avg(list)` | Среднее арифметическое |
| `random` | `random()` | Случайное `float` в `[0, 1)` |
| `randomInt` | `randomInt(lo, hi)` | Случайное `int` в `[lo, hi)` |
| `gcd` | `gcd(a, b)` | Наибольший общий делитель |
| `lcm` | `lcm(a, b)` | Наименьшее общее кратное |

**Примеры:**

```
round(3.14159, 2)  // 3.14
clamp(150, 0, 100) // 100
min(5, 3, 8, 1)    // 1
sum([1, 2, 3, 4])  // 10
avg([10, 20, 30])  // 20.0
log(8, 2)          // 3.0
gcd(12, 8)         // 4
```

---

### Списки

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `append` | `append(list, item)` | Добавить элемент в конец; возвращает новый список |
| `prepend` | `prepend(list, item)` | Добавить элемент в начало; возвращает новый список |
| `concat` | `concat(list1, list2, …)` | Объединить списки |
| `first` | `first(list)` / `first(list, n)` | Первый элемент или первые `n` элементов |
| `last` | `last(list)` / `last(list, n)` | Последний элемент или последние `n` элементов |
| `nth` | `nth(list, i)` | Элемент по индексу `i`; отрицательный — с конца |
| `slice` | `slice(list, start)` / `slice(list, start, end)` | Подсписок; отрицательные индексы считаются с конца |
| `take` | `take(list, n)` | Первые `n` элементов |
| `drop` | `drop(list, n)` | Список без первых `n` элементов |
| `contains` | `contains(list, item)` | Содержит ли список элемент; возвращает `bool` |
| `indexOf` | `indexOf(list, item)` | Первый индекс элемента или `-1` |
| `lastIndexOf` | `lastIndexOf(list, item)` | Последний индекс элемента или `-1` |
| `flatten` | `flatten(list)` | Разгладить на один уровень |
| `flattenAll` | `flattenAll(list)` | Разгладить рекурсивно |
| `unique` | `unique(list)` | Удалить дубликаты (порядок сохраняется) |
| `sort` | `sort(list)` | Сортировка по возрастанию |
| `sortDesc` | `sortDesc(list)` | Сортировка по убыванию |
| `range` | `range(end)` / `range(start, end)` / `range(start, end, step)` | Сгенерировать последовательность целых чисел |
| `chunk` | `chunk(list, size)` | Разбить на подсписки заданного размера |
| `zip` | `zip(list1, list2)` | Попарно объединить: `[[a0,b0], [a1,b1], …]` |
| `without` | `without(list, item1, item2, …)` | Список без указанных элементов |
| `count` | `count(list)` | Количество элементов (псевдоним `len`) |
| `reverse` | `reverse(list)` | Перевернуть список |
| `len` | `len(list)` | Количество элементов |

**Примеры:**

```
append([1, 2], 3)               // [1, 2, 3]
concat([1, 2], [3, 4])          // [1, 2, 3, 4]
first([10, 20, 30], 2)          // [10, 20]
slice([0, 1, 2, 3, 4], 1, 4)   // [1, 2, 3]
range(1, 6)                     // [1, 2, 3, 4, 5]
range(0, 10, 2)                 // [0, 2, 4, 6, 8]
unique([1, 2, 2, 3, 1])         // [1, 2, 3]
sort([3, 1, 4, 1, 5])           // [1, 1, 3, 4, 5]
chunk([1, 2, 3, 4, 5], 2)       // [[1, 2], [3, 4], [5]]
zip([1, 2, 3], ["a", "b", "c"]) // [[1, "a"], [2, "b"], [3, "c"]]
without([1, 2, 3, 4], 2, 4)     // [1, 3]
flatten([[1, 2], [3, [4, 5]]])   // [1, 2, 3, [4, 5]]
flattenAll([[1, [2, [3]]]])      // [1, 2, 3]
```

---

### Словари

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `keys` | `keys(map)` | Список ключей (в порядке добавления) |
| `values` | `values(map)` | Список значений (в порядке добавления) |
| `hasKey` | `hasKey(map, key)` | Есть ли ключ; возвращает `bool` |
| `get` | `get(map, key)` / `get(map, key, default)` | Получить значение или `default`, если ключ отсутствует |
| `set` | `set(map, key, value)` | Вернуть новый словарь с установленным ключом |
| `delete` | `delete(map, key1, key2, …)` | Вернуть новый словарь без указанных ключей |
| `merge` | `merge(map1, map2, …)` | Объединить словари; при совпадении ключей побеждает последний |
| `entries` | `entries(map)` | Список пар `[ключ, значение]` |
| `fromEntries` | `fromEntries([[key, value], …])` | Построить словарь из пар |
| `pick` | `pick(map, key1, key2, …)` | Новый словарь только с указанными ключами |
| `omit` | `omit(map, key1, key2, …)` | Новый словарь без указанных ключей |
| `len` | `len(map)` | Количество ключей |

**Примеры:**

```
keys({"a": 1, "b": 2})                     // ["a", "b"]
get({"x": 10}, "x")                        // 10
get({"x": 10}, "y", 0)                     // 0
set({"a": 1}, "b", 2)                      // {"a": 1, "b": 2}
delete({"a": 1, "b": 2, "c": 3}, "b")     // {"a": 1, "c": 3}
merge({"a": 1}, {"b": 2, "a": 99})        // {"a": 99, "b": 2}
pick({"a": 1, "b": 2, "c": 3}, "a", "c")  // {"a": 1, "c": 3}
omit({"a": 1, "b": 2, "c": 3}, "b")       // {"a": 1, "c": 3}
entries({"x": 1, "y": 2})                 // [["x", 1], ["y", 2]]
fromEntries([["name", "Alice"], ["age", 30]]) // {"name": "Alice", "age": 30}
```

---

### Логика

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `not` | `not(v)` | Логическое отрицание |
| `and` | `and(a, b, …)` | `true`, если все аргументы истинны |
| `or` | `or(a, b, …)` | `true`, если хотя бы один аргумент истинен |
| `xor` | `xor(a, b)` | `true`, если ровно один из двух истинен |
| `if` | `if(cond, then)` / `if(cond, then, else)` | Условие (оба ветки вычисляются заранее) |
| `coalesce` | `coalesce(a, b, …)` | Первое не-`null` значение |
| `defaultTo` | `defaultTo(v, default)` | Возвращает `v`, если оно не `null`, иначе `default` |
| `isTruthy` | `isTruthy(v)` | Является ли значение истинным? |
| `isFalsy` | `isFalsy(v)` | Является ли значение ложным? |

Ложные значения: `false`, `0`, `0.0`, `""`, `"false"`, `"0"`, `null`, `[]`, `{}`. Всё остальное — истинно.

**Примеры:**

```
and(true, 1, "yes")         // true
or(false, 0, "hello")       // true
coalesce(null, null, 42)    // 42
defaultTo(null, "fallback") // "fallback"
if(10 > 5, "big", "small")  // "big"
```

---

### Преобразование и проверка типов

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `toString` | `toString(v)` | Преобразовать любое значение в строку |
| `toInt` | `toInt(v)` | Преобразовать в `int` (float усекается, строка разбирается) |
| `toFloat` | `toFloat(v)` | Преобразовать в `float` |
| `toBool` | `toBool(v)` | Преобразовать в `bool` по правилам истинности |
| `toList` | `toList(v)` | Обернуть не-список в список; `null` → `[]` |
| `typeOf` | `typeOf(v)` | Вернуть название типа: `"null"`, `"bool"`, `"int"`, `"float"`, `"string"`, `"list"`, `"map"` |
| `isNull` | `isNull(v)` | Является ли значение `null`? |
| `isBool` | `isBool(v)` | Является ли значение `bool`? |
| `isInt` | `isInt(v)` | Является ли значение `int`? |
| `isFloat` | `isFloat(v)` | Является ли значение `float`? |
| `isNumber` | `isNumber(v)` | Является ли значение числом (`int` или `float`)? |
| `isString` | `isString(v)` | Является ли значение строкой? |
| `isList` | `isList(v)` | Является ли значение списком? |
| `isMap` | `isMap(v)` | Является ли значение словарём? |

**Примеры:**

```
typeOf(42)        // "int"
typeOf([1, 2])    // "list"
toInt("123")      // 123
toFloat("3.14")   // 3.14
toString(true)    // "true"
toList(null)      // []
toList(5)         // [5]
isNumber(3.14)    // true
```

---

### Прочее

| Функция | Сигнатура | Описание |
|---------|-----------|----------|
| `print` | `print(v, …)` | Вывести аргументы в stdout; возвращает `null` |
| `error` | `error(msg)` | Вернуть ошибку с указанным сообщением |
| `uuid` | `uuid()` | Сгенерировать случайный UUID v4 |
| `now` | `now()` | Текущий Unix-timestamp в секундах (`int`) |
| `nowMs` | `nowMs()` | Текущий Unix-timestamp в миллисекундах (`int`) |
| `dateFormat` | `dateFormat(unixSec, layout)` | Форматировать Unix-timestamp по шаблону времени Go |
| `dateParse` | `dateParse(layout, str)` | Разобрать строку с датой, вернуть Unix-timestamp (`int`) |
| `sleep` | `sleep(ms)` | Приостановить выполнение на `ms` миллисекунд; возвращает `null` |

**Примеры:**

```
uuid()                                 // "f47ac10b-58cc-4372-a567-0e02b2c3d479"
now()                                  // 1711612800
dateFormat(now(), "2006-01-02")        // "2024-03-28"
dateParse("2006-01-02", "2024-01-15") // 1705276800
error("что-то пошло не так")          // возвращает ошибку
```

> `dateFormat` и `dateParse` используют эталонное время Go: `Mon Jan 2 15:04:05 MST 2006`.

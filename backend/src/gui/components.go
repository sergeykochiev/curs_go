package gui

import (
	"database/sql"
	"fmt"

	. "github.com/sergeykochiev/curs/backend/entity"
	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func LabelComponent(children Node, label string) Node {
	return Label(
		Class("flex flex-col p-[8px] bg-gray-50 gap-[4px] border border-[1px] border-solid border-gray-200 font-medium text-[12px] w-full"),
		Text(label),
		children,
	)
}

func InputComponent(t string, ph string, name string, label string, default_value string, required bool) Node {
	return LabelComponent(
		Input(
			Class("px-[12px] w-full py-[6px] text-[16px] font-normal outline-gray-300 focus:outline-gray-400 transition-all hover:bg-gray-200 bg-gray-100 focus:bg-gray-50 outline outline-[1.5px]"),
			Type(t),
			Placeholder(ph),
			Name(name),
			If(required, Required()),
			Value(default_value),
		), label,
	)
}

func SelectComponent[T interface {
	HtmlTemplater
	Identifier
}](arr []T, ph string, getText func(T) string, label string, name string, required bool, id int) Node {
	return LabelComponent(
		Select(
			If(required, Required()),
			Class("w-full px-[12px] w-full font-normal text-[16px] py-[6px] outline-gray-200 focus:outline-gray-400 transition-all hover:bg-gray-200 bg-gray-100 focus:bg-gray-50 outline outline-[1.5px]"),
			Placeholder(ph),
			Name(name),
			Map(arr, func(ent T) Node {
				If(id == ent.GetId(), Selected())
				return Option(Text(getText(ent)), Value(fmt.Sprintf("%d", ent.GetId())))
			}),
		), label,
	)
}

func TailwindScript() Node {
	return Script(
		Src("https://unpkg.com/@tailwindcss/browser@4"),
	)
}

func DataTableComponent[T interface {
	HtmlTemplater
	Identifier
}](ent T, arr []T) Node {
	return Div(
		Class("flex flex-col w-full gap-[8px]"),
		Div(
			Class("flex outline-gray-200 outline outline-[1.5px]"),
			ent.GetTableHeader(),
		),
		If(len(arr) > 0, Div(
			Class("flex flex-col gap-[2px]"),
			Map(arr, func(ent T) Node {
				return Div(
					Class("bg-gray-100 flex"),
					ent.ToHtmlDataRow(),
				)
			}),
		)),
	)
}

func RootComponent(children Node) Node {
	return HTML(
		Head(TailwindScript()),
		Body(
			Class("flex justify-center"),
			children,
		),
	)
}

func MainPageButtonComponent(href string, text string) Node {
	return A(
		Class("self-end w-full py-[16px] px-[16px] grid place-items-center font-medium text-slate-700 text-[16px] outline-gray-300 hover:outline-gray-400 bg-gray-100 hover:bg-gray-200 transition-all cursor-pointer active:scale-[0.95] outline-[1.5px]"),
		Text(text),
		Href(href),
	)
}

func MainPageSectionComponent(heading string, children Group) Node {
	return Section(
		Class("flex flex-col w-full gap-[6px]"),
		H2(
			Text(heading),
			Class("text-[20px] font-semibold"),
		),
		Div(
			Class("grid grid-cols-2 gap-[8px] w-full"),
			children,
		),
	)
}

func MainPageComponent() Node {
	return RootComponent(
		Main(
			MainWrapperClass(),
			MainPageSectionComponent("Функции менеджмента", Group{
				MainPageButtonComponent("/create_order", "Создать заказ"),
				MainPageButtonComponent("/end_order", "Завершить заказ"),
				MainPageButtonComponent("/create_spending", "Завести трату ресурса"),
				MainPageButtonComponent("/create_resupply", "Завести поставку ресурса"),
			}),
			MainPageSectionComponent("Просмотр данных", Group{
				MainPageButtonComponent("/resources", "Посмотреть все ресурсы на складе"),
				MainPageButtonComponent("/orders", "Посмотреть все заказы"),
				MainPageButtonComponent("/resource_resupplies", "Посмотреть все поставки ресурсов"),
				MainPageButtonComponent("/resource_spendings", "Посмотреть все траты ресурсов"),
			}),
			MainPageSectionComponent("Продвинутые функции", Group{
				MainPageButtonComponent("/create_resource", "Добавить ресурс на склад вручную"),
			}),
		),
	)
}

func MainWrapperClass() Node {
	return Class("flex flex-col mt-[30px] max-w-[960px] gap-[16px] grid grid-cols-1 w-full")
}

func DataPageComponent[T interface {
	HtmlTemplater
	Identifier
}](ent T, arr []T, db *sql.DB) Node {
	return RootComponent(
		Main(
			MainWrapperClass(),
			H1(
				Class("text-[20px] font-semibold"),
				Text(ent.GetReadableName()),
			),
			Class("flex flex-col gap-[12px]"),
			DataTableComponent(ent, arr),
		),
	)
}

func ButtonComponent(text string) Node {
	return Button(
		Class("self-end px-[16px] py-[6px] font-medium text-[14px] outline-gray-400 bg-gray-100 hover:bg-gray-200 transition-all cursor-pointer active:scale-[0.95] outline-[1.5px]"),
		Text(text),
	)
}

func UserFormComponent(signup bool) Node {
	return RootComponent(
		Main(
			MainWrapperClass(),
			Form(
				Method("post"),
				Class("flex flex-col gap-[12px]"),
				H2(Text(ConditionalArg(signup, "Регистрация", "Вход"))),
				InputComponent("text", "Ivan2000Rus", "name", "Имя пользователя", "", true),
				InputComponent("password", "Не менее 8-ми символов", "password", "Пароль", "", true),
				If(signup, InputComponent("password", "Должен совпадать с паролем выше", "repeat_password", "Повторите пароль", "", true)),
				A(
					Href(ConditionalArg(signup, "/login", "/signup")),
					Text(ConditionalArg(signup, "Есть аккаунт? Войти", "Нет аккаунта? Зарегистрироваться")),
				),
				ButtonComponent(ConditionalArg(signup, "Зарегистрироваться", "Войти")),
			),
		),
	)
}

func CreateOrderFormComponent() Node {
	return RootComponent(
		Main(
			MainWrapperClass(),
			Form(
				Method("post"),
				Class("flex flex-col gap-[12px]"),
				H2(
					Text("Создать заказ"),
					Class("text-[20px] font-semibold"),
				),
				InputComponent("text", "", "name", "Название заказа", "", true),
				InputComponent("text", "", "client_name", "Имя клиента", "", true),
				InputComponent("number", "", "client_phone", "Телефон клиента", "", true),
				InputComponent("date", "", "date_created", "Дата создания", "", true),
				ButtonComponent("Создать"),
			),
		),
	)
}

func EndOrderComponent(arr []*OrderEntity) Node {
	return RootComponent(
		Main(
			MainWrapperClass(),
			Form(
				Method("post"),
				Class("flex flex-col gap-[12px]"),
				H2(
					Text("Завершить заказ"),
					Class("text-[20px] font-semibold"),
				),
				SelectComponent(arr, "Выберите необходимый заказ", func(ent *OrderEntity) string { return ent.Name }, "Заказ"+func() string {
					if len(arr) == 0 {
						return " (нет подходящих)"
					} else {
						return ""
					}
				}(), "id", true, -1),
				InputComponent("date", "", "date_ended", "Дата выполнения", "", true),
				ButtonComponent("Завершить"),
			),
		),
	)
}

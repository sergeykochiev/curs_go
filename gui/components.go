package gui

import (
	"fmt"

	"github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	icons "maragu.dev/gomponents-heroicons/v3/outline"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func NotFoundPage() Node {
	return Div(
		Class("w-screen h-screen grid place-items-center text-[32px] font-bold"),
		Text("404"),
	)
}

func He1(text string) Node {
	return H1(
		Class("font-bold text-[24px]"),
		Text(text),
	)
}

func He2(text string) Node {
	return H2(
		Class("font-semibold text-[20px]"),
		Text(text),
	)
}

func EntityListPageVerticalLayout[T interface {
	types.HtmlTemplater
	types.Identifier
}](ent T, arr []T) Node {
	return Div(
		Class("w-full flex gap-[24px]"),
		FiltersPanelComponent(ent),
		DataTableComponent(ent, arr),
	)
}

func EntityListPage[T interface {
	types.HtmlTemplater
	types.Identifier
}](ent T, arr []T) Node {
	return PageComponent(EntityListPageVerticalLayout(ent, arr), ent.GetReadableName(), "На главную", "/", ButtonComponent("Создать", A, Href(ent.TableName()+"/create")))
}

func EntityPage[T interface {
	types.HtmlTemplater
	types.Identifier
}](ent T) Node {
	return PageComponent(ent.GetEntityPage(true), fmt.Sprintf("%s #%d", ent.GetReadableName(), ent.GetId()), "К таблице", "/"+ent.TableName(), ent.GetEntityPageButtons())
}

func PageComponent(content Node, heading string, link_text string, link_href string, buttons ...Node) Node {
	return RootComponent(
		Main(
			MainWrapperClass(""),
			Div(
				Class("w-max"),
				A(
					Class("transition-all text-[14px] flex items-center after:content-[''] after:transition-all after:w-0 after:z-[-1] text-gray-800 gap-[4px] after:bg-gray-200 after:h-full hover:after:w-full relative after:absolute after:bottom-0 after:left-0"),
					icons.ArrowLeft(Class("h-4 w-4")),
					Text(link_text),
					Href(link_href),
				),
			),
			Div(
				Class("flex items-center justify-between font-semibold"),
				He1(heading),
				Div(
					Class("flex gap-[8px] items-center"),
					Map(buttons, func(b Node) Node { return b }),
				),
			),
			Class("flex flex-col gap-[12px] w-full"),
			content,
		),
	)
}

func RelationCardComponent[T interface {
	types.HtmlTemplater
	types.Identifier
}](heading string, ent T) Node {
	return RelationCardCoreComponent(
		heading, GetOneHref(ent), ent.GetEntityPage(false),
	)
}

func RelationCardCoreComponent(heading string, href string, children Group) Node {
	return A(
		Href(href),
		Class("transition-all bg-gray-100 flex flex-col gap-[8px] p-[8px] hover:bg-gray-200 outline outline-[1.5px] outline-gray-400"),
		He2(heading),
		children,
	)
}

func RelationCardArrComponent[T interface {
	types.HtmlTemplater
	types.Identifier
}](heading string, arr []T, f func(ent T) Node) Node {
	return MainDataContainerComponent(Div, Group{
		He2(heading),
		If(len(arr) > 0, Map(arr, func(ent T) Node {
			return f(ent)
		})),
		If(len(arr) == 0, Div(Class("grid place-items-center w-full"), Text("Нет данных"))),
	}, true)
}

func LabeledFieldComponent(label string, value string) Node {
	return Div(
		Class("flex items-center justify-between w-full"),
		P(
			Class("font-medium"),
			Text(label),
		),
		Text(value),
	)
}

//	func LabelComponent(children Node, label string) Node {
//		return Label(
//			Class("flex flex-col p-[8px] bg-gray-50 gap-[4px] border border-[1px] border-solid border-gray-200 font-medium text-[12px] w-full"),
//			Text(label),
//			children,
//		)
//	}

func LabelComponent(children Node, label string) Node {
	return Label(
		Class("flex flex-col gap-[4px] font-medium text-[12px] w-full"),
		Text(label),
		children,
	)
}

func InputComponent(t string, ph string, name string, default_value string, required bool) Node {
	return Input(
		Class("px-[12px] w-full py-[6px] text-[16px] font-normal outline-gray-300 focus:outline-gray-400 transition-all hover:bg-gray-200 bg-gray-100 focus:bg-gray-50 outline outline-[1.5px]"),
		Type(t),
		Placeholder(ph),
		Name(name),
		If(required, Required()),
		Value(default_value),
	)
}

func LabeledInputComponent(t string, ph string, name string, label string, default_value string, required bool) Node {
	return LabelComponent(
		InputComponent(t, ph, name, default_value, required), label,
	)
}

func SelectComponent[T interface {
	types.HtmlTemplater
	types.Identifier
}](arr []T, ph string, getText func(T) string, label string, name string, required bool, id int64) Node {
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
		Src("/tailwind.js"),
	)
}

func MainDataContainerComponent(as func(...Node) Node, children Group, full_width bool) Node {
	return as(
		Class("shadow-sm p-[12px] outline-gray-200 outline outline-[1px] flex flex-col gap-[12px]"+ConditionalArg(full_width, " w-full", "")),
		children,
	)
}

func FiltersPanelComponent[T interface {
	types.HtmlTemplater
	types.Identifier
}](ent T) Node {
	return MainDataContainerComponent(Form, Group{
		He2("Фильтры"),
		ent.GetFilters(),
		ButtonComponent("Применить", Button),
	}, false)
}

func DataTableComponent[T interface {
	types.HtmlTemplater
	types.Identifier
}](ent T, arr []T) Node {
	return MainDataContainerComponent(Div, Group{
		He2("Данные"),
		Table(
			Class("border-collapse"),
			THead(
				Tr(
					Class("outline-gray-200 whitespace-nowrap outline *:py-[2px] *:text-center *:border-r *:last:border-none *:border-gray-200"),
					ent.GetTableHeader(),
				),
			),
			If(len(arr) > 0, TBody(
				Map(arr, func(ent T) Node {
					return Tr(
						Class("transition-all relative bg-gray-100 hover:bg-gray-200 outline outline-gray-400 *:py-[2px] *:text-center"),
						ent.GetDataRow(),
					)
				}),
			)),
		),
	}, true)
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

func HeadingNavTabComponent(label string, name string, id string, default_checked bool) Node {
	return Label(
		Class("group/tab has-[input:checked]:cursor-default cursor-pointer"),
		H1(
			Text(label),
			Class("group-has-[input:checked]/tab:text-gray-800 group-[&:not(:has(input:checked))]:underline text-gray-400 transition-all font-bold text-[24px]"),
		),
		Input(
			Type("radio"),
			Name(name),
			ID(id),
			If(default_checked, Checked()),
			Class("absolute hidden"),
		),
	)
}

func MainPageComponent() Node {
	return RootComponent(
		Main(
			MainWrapperClass("group/panel"),
			Div(
				Class("flex gap-[12px]"),
				HeadingNavTabComponent("Обычный режим", "input-toggle-mode", "input-toggle-basic", true),
				HeadingNavTabComponent("Экспертный режим", "input-toggle-mode", "input-toggle-expert", false),
			),
			Div(
				Class("group-has-[input#input-toggle-basic:checked]/panel:hidden grid grid-cols-2 gap-[8px] w-full"),
				MainPageButtonComponent("/resource", "Ресурсы на складе"),
				MainPageButtonComponent("/order", "Заказы"),
				MainPageButtonComponent("/resource_resupply", "Поставки ресурсов"),
				MainPageButtonComponent("/resource_spending", "Траты ресурсов на заказ"),
				MainPageButtonComponent("/item", "Товары"),
				MainPageButtonComponent("/order_item_fulfillment", "Предоставления товаров в рамках заказа"),
				MainPageButtonComponent("/item_resource_need", "Необходимость ресурса на товар"),
			),
			Div(
				Class("group-has-[input#input-toggle-expert:checked]/panel:hidden grid grid-cols-2 gap-[8px] w-full"),
				MainPageButtonComponent("/resource_resupply/create", "Добавить поставку ресурса"),
				MainPageButtonComponent("/order", "Посмотреть заказы"),
			),
		),
	)
}

func MainWrapperClass(class string) Node {
	return Class("flex flex-col mt-[30px] max-w-[1440px] gap-[16px] grid grid-cols-1 w-full pb-[32px] " + class)
}

func ButtonComponent(text string, as func(children ...Node) Node, children ...Node) Node {
	return as(
		Class("self-end px-[16px] py-[6px] font-medium text-[14px] outline-gray-400 bg-gray-100 hover:bg-gray-200 transition-all cursor-pointer active:scale-[0.95] outline-[1.5px]"),
		Text(text),
		Map(children, func(c Node) Node { return c }),
	)
}

func UserFormComponent(signup bool) Node {
	return RootComponent(
		Main(
			MainWrapperClass(""),
			Form(
				Method("post"),
				Class("flex flex-col gap-[12px]"),
				He1(ConditionalArg(signup, "Регистрация", "Вход")),
				LabeledInputComponent("text", "Ivan2000Rus", "name", "Имя пользователя", "", true),
				LabeledInputComponent("password", "Не менее 8-ми символов", "password", "Пароль", "", true),
				If(signup, LabeledInputComponent("password", "Должен совпадать с паролем выше", "repeat_password", "Повторите пароль", "", true)),
				A(
					Href(ConditionalArg(signup, "/login", "/signup")),
					Text(ConditionalArg(signup, "Есть аккаунт? Войти", "Нет аккаунта? Зарегистрироваться")),
				),
				ButtonComponent(ConditionalArg(signup, "Зарегистрироваться", "Войти"), Button),
			),
		),
	)
}

func CreateFormComponent(name string, fields Group) Node {
	return RootComponent(
		Main(
			MainWrapperClass(""),
			Form(
				Method("post"),
				Class("flex flex-col gap-[12px]"),
				He2(fmt.Sprintf("Создать %s", name)),
				fields,
				ButtonComponent("Создать", Button),
			),
		),
	)
}

func TableDataComponent(value string, as func(children ...Node) Node, href string) Node {
	return as(
		Text(value),
		If(href != "", A(
			Href(href),
			Class("left-0 top-0 absolute w-full h-full z-1 bg-transparent cursor-pointer"),
		)))
}

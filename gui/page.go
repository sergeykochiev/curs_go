package gui

import (
	"fmt"

	"github.com/sergeykochiev/curs/backend/types"
	util "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func MainPage() Node {
	return RootComponent(
		Main(
			MainWrapperClass("group/panel"),
			Div(
				Class("flex gap-[12px]"),
				HeadingNavTabComponent("Обычный режим", "input-toggle-mode", "input-toggle-basic", true),
				HeadingNavTabComponent("Экспертный режим", "input-toggle-mode", "input-toggle-expert", false),
			),
			Div(
				Class("group-has-[input#input-toggle-basic:checked]/panel:hidden"),
				DoubleGridComponent(Group{
					MainPageButtonComponent("/resource", "Ресурсы на складе"),
					MainPageButtonComponent("/order", "Заказы"),
					MainPageButtonComponent("/resource_resupply", "Поставки ресурсов"),
					MainPageButtonComponent("/resource_spending", "Траты ресурсов на заказ"),
					MainPageButtonComponent("/item", "Товары"),
					MainPageButtonComponent("/order_item_fulfillment", "Предоставления товаров в рамках заказа"),
					MainPageButtonComponent("/item_resource_need", "Необходимость ресурса на товар"),
				}),
			),
			Div(
				Class("group-has-[input#input-toggle-expert:checked]/panel:hidden"),
				DoubleGridComponent(Group{
					MainPageButtonComponent("/resource_resupply/create", "Добавить поставку ресурса"),
					MainPageButtonComponent("/resource", "Посмотреть ресурсы на складе"),
					MainPageButtonComponent("/order", "Посмотреть заказы"),
					MainPageButtonComponent("/item_popularity", "Получить отчет о популярности товаров/услуг"),
					MainPageButtonComponent("/resource_spendings", "Получить отчет о тратах ресурсов"),
				}),
			),
		),
	)
}

func DatedReportFormPage(heading string) Node {
	return FormPageComponent(Group{
		Method("POST"),
		He2(heading),
		LabeledInputComponent("date", "", "date_lo", "С (дата)", "", false),
		LabeledInputComponent("date", "", "date_hi", "По (дата)", "", false),
		ButtonComponent("Создать отчет", Button),
	})
}

func UserFormPage(signup bool) Node {
	return FormPageComponent(Group{
		Method("POST"),
		He1(util.ConditionalArg(signup, "Регистрация", "Вход")),
		LabeledInputComponent("text", "Иванов Иван", "name", "Имя пользователя", "", true),
		LabeledInputComponent("password", "Не менее 8-ми символов", "password", "Пароль", "", true),
		If(signup, LabeledInputComponent("password", "Должен совпадать с паролем выше", "repeat_password", "Повторите пароль", "", true)),
		A(
			Href(util.ConditionalArg(signup, "/login", "/signup")),
			Text(util.ConditionalArg(signup, "Есть аккаунт? Войти", "Нет аккаунта? Зарегистрироваться")),
		),
		ButtonComponent(util.ConditionalArg(signup, "Зарегистрироваться", "Войти"), Button),
	})
}

func NotFoundPage() Node {
	return Div(
		Class("w-screen h-screen grid place-items-center text-[32px] font-bold"),
		Text("404"),
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

func CreateFormPage(name string, fields Group) Node {
	return FormPageComponent(Group{
		Method("POST"),
		He2(fmt.Sprintf("Создать %s", name)),
		fields,
		ButtonComponent("Создать", Button),
	})
}

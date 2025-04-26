package gui

import (
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func DateFilterComponent(label string, name string) Node {
	return LabelComponent(Div(
		Class("flex items-center gap-[4px]"),
		InputComponent("date", "", name+"_lo", "", false),
		Text(":"),
		InputComponent("date", "", name+"_hi", "", false),
	), label)
}

func BoolFilterComponent(label string, name string) Node {
	return LabelComponent(Div(
		Class("flex items-center justify-between"),
		Label(
			Class("cursor-pointer flex gap-[4px] w-full items-center"),
			Input(Type("radio"), Value("1"), Name(name)),
			Text("Да"),
		),
		Label(
			Class("cursor-pointer flex gap-[4px] w-full items-center"),
			Input(Type("radio"), Value("0"), Name(name)),
			Text("Нет"),
		),
	), label)
}

func StringFilterComponent(label string, name string) Node {
	return LabeledInputComponent("text", "", name, label, "", false)
}

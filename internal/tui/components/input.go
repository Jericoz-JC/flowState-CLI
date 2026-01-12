package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	textinput textinput.Model
	focused   bool
}

func NewTextInput(placeholder string) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.Prompt = "> "
	ti.CharLimit = 500

	return TextInputModel{textinput: ti, focused: true}
}

func (m *TextInputModel) SetValue(value string) {
	m.textinput.SetValue(value)
}

func (m *TextInputModel) Value() string {
	return m.textinput.Value()
}

func (m *TextInputModel) Focus() {
	m.focused = true
	m.textinput.Focus()
}

func (m *TextInputModel) Blur() {
	m.focused = false
	m.textinput.Blur()
}

func (m *TextInputModel) Focused() bool {
	return m.focused
}

func (m *TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	ti, cmd := m.textinput.Update(msg)
	return TextInputModel{textinput: ti, focused: m.focused}, cmd
}

func (m *TextInputModel) View() string {
	return m.textinput.View()
}

type TextAreaModel struct {
	textarea textarea.Model
	focused  bool
}

func NewTextArea(placeholder string) TextAreaModel {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.Focus()
	ta.Prompt = "| "
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false
	ta.SetHeight(10)

	return TextAreaModel{textarea: ta, focused: true}
}

func (m *TextAreaModel) SetValue(value string) {
	m.textarea.SetValue(value)
}

func (m *TextAreaModel) Value() string {
	return m.textarea.Value()
}

func (m *TextAreaModel) Focus() {
	m.focused = true
	m.textarea.Focus()
}

func (m *TextAreaModel) Blur() {
	m.focused = false
	m.textarea.Blur()
}

func (m *TextAreaModel) Focused() bool {
	return m.focused
}

func (m *TextAreaModel) Update(msg tea.Msg) (TextAreaModel, tea.Cmd) {
	ta, cmd := m.textarea.Update(msg)
	return TextAreaModel{textarea: ta, focused: m.focused}, cmd
}

func (m *TextAreaModel) View() string {
	return m.textarea.View()
}

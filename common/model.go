package common

type Recipient struct {
	Email string
	Name string
	Type string
}

type Mail struct {
	Recipients []Recipient
	FromEmail string
	FromName string
	Subject string
	HTML string
	Text string
}
func (m *Mail) AddRecipient(email, name, t string) {
	recipient := Recipient{Email:email, Name:name, Type:t}
	m.Recipients = append(m.Recipients, recipient)
}

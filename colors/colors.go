package colors

import "github.com/muesli/termenv"

var (
	color          = termenv.EnvColorProfile().Color
	ErrorKeyword   = termenv.Style{}.Foreground(color("#E06C75")).Styled
	SuccessKeyword = termenv.Style{}.Foreground(color("#98C379")).Styled
	InfoKeyword    = termenv.Style{}.Foreground(color("#61AFEF")).Styled
	Heading        = termenv.Style{}.Foreground(color("#61AFEF")).Styled
)

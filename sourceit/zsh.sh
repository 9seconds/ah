#!/usr/bin/env zsh
# vim: set noexpandtab shiftwidth=4:

__auto_ah() {
	BUFFER=$(ah at "${BUFFER}")
	zle accept-line
}

zle -N __auto_ah_widget __auto_ah
bindkey '^J' __auto_ah_widget
bindkey '^M' __auto_ah_widget

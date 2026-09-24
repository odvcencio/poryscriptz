; poryz syntax captures for the embedded gotreesitter grammar.
(comment) @comment
(interpreted_string_literal) @string
(raw_string_literal) @string
(int_literal) @number

(script_declaration name: (identifier) @type)
(label_declaration name: (identifier) @label)
(movement_declaration name: (identifier) @label)
(call_expression function: (identifier) @function)

(identifier) @constant (#match? @constant "^(FLAG|VAR|ITEM|SPECIES|MOVE|MAP)_[A-Za-z0-9_]+$")

"script" @keyword
"label" @keyword
"movement" @keyword
"include" @keyword
"entries" @keyword
"header" @keyword
"init" @keyword
"source" @keyword
"dialect" @keyword
"pret" @keyword
"if" @keyword
"else" @keyword
"for" @keyword
"switch" @keyword
"case" @keyword
"default" @keyword

import re
import os

f = open("./tex/notes.MD")
a = open("./tex/appendix.md")

text = f.read()
appendix = a.read()

r = re.compile(r"(\$\$)\s(\\begin.*?)(\$\$)", re.DOTALL)

out = open("./out/out.MD", "w+")
out_appendix = open("./out/appendix.MD", "w+")

text = re.sub(r, r"\2", text)
appendix = re.sub(r, r"\2", appendix)

out.write(text)
out_appendix.write(appendix)

out.close()
out_appendix.close()

code = os.system(
    "pandoc --verbose -H ./tex/header.tex --template ./tex/template.tex --number-sections --toc -V geometry:margin=1in --bibliography=./tex/lit.bib --csl=./tex/ieee.csl ./out/out.MD ./out/appendix.MD ./tex/references.md -o notes.pdf")
exit(code)

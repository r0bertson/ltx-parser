# ltx-parser
A parser for LINDO Text file.

This repo is a helper that converts a LP on LINDO's model sintax into a struct.

#####LINDO'S MODEL SINTAX:
```
MAX 10 X1 + 15 X2
SUBJECT TO
 X1 < 10
 X2 < 12
 X1 + 2 X1 < 16
END
```

##### TODO:
1. Variable and constraint names must be at most 8 characters
2. Improve test cases
3. Add constraint naming support.
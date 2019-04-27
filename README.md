# ltx-parser #

This repo is a helper that converts a linear problem written with LINDO's model syntax (```.ltx``` files) into a golang struct.
All of LINDO's syntax constraints are implemented in the parser, including:

1. A problem must start with one of the following keywords: ```MAX```, ```MAXIMIZE```, ```MAXIMISE```, ```MIN```, ```MINIMIZE``` or ```MINIMISE```
2. Variables must begin with a letter(a-z A-Z) and have at most eight characters.
3. < > = + - ! ) cannot be used inside variable names.
4. The constraints must be preceeded by one of the following keywords: ```SUBJECT TO```,```SUCH THAT```,```ST``` or ```S.T.```
5. The righthand side of a constraint can only have a constant number
6. The lefthand side of a constraint can only have variables
7. The operators ```<``` or ```>``` will be replaced by ```<=``` or ```>=``` automatically.
7. A problem must end with the keyword ```END```

The parser currently only implements the essential and mandatory features of LINDO's syntax model (objective function and constraints). Optional features, such as FREE and SLB/SUB modelling statements, will be added in the future.

##### Exemple of a linear problem written with LINDO's syntax:
```
MAX 10 X1 + 15 X2
SUBJECT TO
 X1 < 10
 X2 < 12
 X1 + 2 X1 < 16
END
```

### TODO:
1. Convert all input variables to UPPERCASE.
2. Improve test cases.
3. Add constraint naming support.
4. Implement optional features of LINDO's syntax model.

### Need help to understand the source? ###
If you have any question, feel free to get in touch:

* [email@robertsonlima.com](mailto:email@robertsonlima.com)
* [Linkedin](https://www.linkedin.com/in/r0bertson/)
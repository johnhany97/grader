# Automatic Programming Assignments Grader

```
Usage of main.go:
  -file string
        Name of file (required)
  -folder string
        Folder in which the file is stored (if not in same dir as main.go)
  -out string
        Where to store the results of the grader (required)
  -scheme string
        Marking scheme to follow when grading the assignment (required)
```

## Installation
The following Docker containers are required for the app to support the most famous languages; Java, C++ & Python. The list of supported languages is to be added below.
```
docker pull dexec/lang-java
docker pull dexec/cpp
docker pull dexec/python
```

## Example usage
```
go run main.go -file="Solution.java" -folder="examples"
```

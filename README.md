# Automatic Programming Assignments Grader

```
Usage of main.go:
  -schema string
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
go run main.go -schema="examples/inputOutput/javaSchema.json"
```

# Assignment 1: Github project data
## Instructions:

Develop a service that will consume a given GitHub project URI and will return 
the associated user account/organisation, an indication of a programming
language(s) used, and the account name of the top committer, that is, the 
contributor with the largest number of commits to the project.

## Service Specification:

The service has to be deployed on either Google Compute Engine or Heroku and expose an API that commits to the following specifications. The service has to be written in Go programming language, must pass Lint and Vet without warnings, and must have at least 20% test coverage. The service is stateless, should not store or record any information, and it should allow concurrent access from multiple clients at the same time. 

## Invocation

### Base path: 
```
GET /projectinfo/v1/[url]
```

### Example: 
```
http://localhost:8080/projectinfo/v1/github.com/apache/kafka
```

### Response payload:

```jsonschema
{

    "project": {
        "type": "string"
    },
    "owner": {
        "type": "string"
    },
    "committer": {
        "type": "string"
    },
    
    "commits": {
        "type": "number"
    },
    "language": {
        "type": "array",
        "items": {
            "type": "string"
        }
    }
}
```

### Example: 
```json
{

    "project": "github.com/apache/kafka",

    "owner": "apache",

    "committer": "enothereska",

    "commits": 19,

    "language": ["Java", "Scala", "Python", "Shell", "Batchfile"]

}
```

## Formal aspects:

This assignment is worth 10% of your total mark. The assignment is individual. Code snippets used from the web and other alternative sources (StackOverflow, tutorials) must be clearly attributed to the original source, in the source code. 

Due date: 21st September 2017

Submission details will be provided closer to the deadline.
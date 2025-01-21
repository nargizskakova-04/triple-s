# triple-s
# Learning Objectives
HTTP
Basic networking concepts
REST API
# Abstract
In this project, you will develop a tool called triple-s, designed to implement a simplified version of S3 (Simple Storage Service) object storage. This tool will provide a REST API that allows clients to interact with the storage system, offering core functionalities such as creating and managing buckets, uploading, retrieving, and deleting files, as well as handling object metadata. The project aims to demonstrate key concepts of RESTful API design, basic networking, and data management, providing a practical foundation for understanding cloud storage solutions.

# Context
Have you ever wondered how cloud storage services, like Amazon S3, manage to store and retrieve files seamlessly over the internet?

Cloud storage systems like S3 provide a highly scalable, reliable, and low-latency way to store data that can be accessed from anywhere. These systems organize data in containers called "buckets" and manage objects (files) within these buckets, allowing for operations like uploading, retrieving, and deleting files, as well as storing metadata about each object.

The triple-s project is a simplified version of S3, designed to give you hands-on experience with the core principles behind such storage systems. Imagine a web service where you can create virtual containers, store files, and access them via a simple URL or API call. While commercial solutions like S3 are highly complex and built for massive scale, triple-s focuses on the essentials:

Buckets: Think of them as folders or containers that hold your files.
Objects: These are the actual files stored within buckets, along with metadata such as size, type, and timestamps.
REST API: A set of HTTP-based operations that allow clients to interact with the storage system, performing actions like creating buckets, uploading files, and retrieving data.
XML Responses: All API responses must be in XML format, in compliance with the Amazon S3 specification.
This project is a practical exploration of how such storage solutions operate under the hood, including how they handle data transfer, manage requests, and more. By building triple-s, you'll dive into key topics like RESTful API design, networking fundamentals, equipping you with the foundational knowledge to understand and even create cloud storage systems.

Whether you're aiming to grasp the basics of cloud storage or prepare for working on real-world distributed systems, triple-s offers a hands-on approach to learning these essential concepts.

# Resources
Read about net/http package here and here
Read about HTTP Server here
Read about HTTP response status codes here
Read about REST and API here, here and here.
Read about 3S API here (pull out only the necessary information)
# General Criteria
Your code MUST be written in accordance with gofumpt. If not, you will be graded 0 automatically.
Your program MUST be able to compile successfully.
Your program MUST not exit unexpectedly (any panics: nil-pointer dereference, index out of range etc.). If so, you will be get 0 during the defence.
Only standard packages are allowed. If not, you will get 0 grade.
The project MUST be compiled by the following command in the project's root directory:
$ go build -o triple-s . 
If an error occurs during startup (e.g., invalid command-line arguments, failure to bind to a port), the program must exit with a non-zero status code and display a clear, understandable error message.
During normal operation, the server must handle errors gracefully, returning appropriate HTTP status codes to the client without crashing.
# Mandatory Part
## Project Initialization
#### Overview:
Before diving into the implementation of RESTful services and storage management, setting up a basic HTTP server is the crucial first step. This foundational setup will act as the backbone for all subsequent features in the triple-s project.

#### Details:
Task: Establish a basic HTTP server using Go’s net/http package.
Objective: To prepare an operational framework that will support the REST API functionalities for managing buckets and objects.
Requirements:
Server Configuration: The server should listen on a configurable port, allowing flexibility in deployment environments.
Request Handling: It must efficiently process incoming HTTP requests and deliver corresponding responses, ensuring initial connectivity and functionality.
Error Handling: Implement error handling to manage and log server errors gracefully. This includes setting up a proper shutdown procedure to handle interruptions or failures smoothly.
Initiating the project with the HTTP server setup ensures that you have a stable and tested platform to build upon. This step is vital for testing the functionality of each API endpoint as they are developed and ensures that the project's infrastructure is robust and reliable from the start.

Your program should take the port number and the path to the directory where the files will be stored as arguments.

## Bucket Management
You will implement core functionalities to manage storage containers, known as "buckets" in the S3 paradigm. Buckets are fundamental units in the storage system where files (objects) are stored. This phase involves creating, listing, and deleting buckets via REST API endpoints, ensuring all responses conform to the XML format required by Amazon S3 specifications.

Note: Authentication and authorization are outside the scope of this project. All operations can be performed without credentials.

### API Endpoints for Bucket Management
#### You will create three primary API endpoints:
#### 1. Create a Bucket:
HTTP Method: PUT
Endpoint: /{BucketName}
Request Body: Empty. Additional parameters can be passed in the request headers.
Behavior:
Validate the bucket name to ensure it meets Amazon S3 naming requirements (3-63 characters, only lowercase letters, numbers, hyphens, and periods).
Ensure the bucket name is unique across the entire storage system.
If the bucket name is valid and unique, create a new entry in the bucket metadata storage.
Return a 200 OK status code and details of the created bucket, or an appropriate error message if the creation fails (e.g., 400 Bad Request for invalid names, 409 Conflict for duplicate names).
Rely on the documentation

#### 2. List All Buckets:
HTTP Method: GET
Endpoint: /
Behavior:
Read the bucket metadata from the storage (e.g., a CSV file).
Return an XML response containing a list of all matching buckets, including metadata like creation time, last modified time, etc.
Respond with a 200 OK status code and the XML list of buckets.
#### 3. Delete a Bucket:
HTTP Method: DELETE
Endpoint: /{BucketName}
Behavior:
Check if the specified bucket exists by looking it up in the bucket metadata storage.
Ensure the bucket is empty (no objects are stored in it) before deletion.
If the bucket exists and is empty, remove it from the metadata storage.
Return a 204 No Content status code if the deletion is successful, or an error message in XML format if the bucket does not exist or is not empty (e.g., 404 Not Found for a non-existent bucket, 409 Conflict for a non-empty bucket).
Don't forget to process the data and save the corresponding metadata in your CSV file.

### Ensuring Unique and Valid Bucket Names:
#### Bucket Naming Conventions:
Bucket names must be unique across the system.
Names should be between 3 and 63 characters long.
Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
Must not be formatted as an IP address (e.g., 192.168.0.1).
Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.
#### Validation Implementation
Use regular expressions to enforce naming rules.
Check the uniqueness of a bucket name by reading the existing entries from the CSV metadata file.
If the bucket name does not meet the rules, return a 400 Bad Request response with a relevant error message.
### Example:
##### Scenario 1: Bucket Creation
A client sends a PUT request to /{BucketName} with the name my-bucket.
The server checks for the validity and uniqueness of the bucket name, then creates an entry in the bucket metadata storage (e.g., buckets.csv).
The server responds with 200 OK and the details of the new bucket or an appropriate error message if the creation fails.
##### Scenario 2: Listing Buckets
A client sends a GET request to /.
The server reads the bucket metadata storage (e.g., buckets.csv) and returns an XML list of all bucket names and metadata.
The server responds with a 200 OK status code.
##### Scenario 3: Deleting a Bucket
A client sends a DELETE request to /{BucketName} for the bucket my-bucket.
The server checks if my-bucket exists and is empty.
If the conditions are met, the bucket is removed from the bucket metadata storage (e.g., buckets.csv).
## Object Operations
This part of the project focuses on implementing the functionality to handle objects (files) stored within buckets. You will create REST API endpoints to upload, retrieve, and delete objects. Each operation will interact with files stored on the disk and update metadata stored in CSV files to keep track of the objects and their attributes.

### Object Key
An object key is a unique identifier for an object (such as a file) stored within a specific bucket in a storage system.

### API Endpoints for Object Operations
You will implement three main API endpoints to handle object operations:

#### 1. Upload a New Object:
HTTP Method: PUT
Endpoint: /{BucketName}/{ObjectKey}
Request Body: Binary data of the object (file content).
Headers:
Content-Type: The object's data type.
Content-Length: The length of the content in bytes.
Behavior:
Verify if the specified bucket {BucketName} exists by reading from the bucket metadata.
Validate the object key {ObjectKey}.
Save the object content to a file in a directory named after the bucket (data/{BucketName}/).
Store object metadata in a CSV file (data/{BucketName}/objects.csv).
Respond with a 200 status code or an appropriate error message if the upload fails.
Note: In this project, if an object with the same name already exists, it must be overwritten.
Check out the examples.

#### 2. Retrieve an Object:
HTTP Method: GET
Endpoint: /{BucketName}/{ObjectKey}
Behavior:
Verify if the bucket {BucketName} exists.
Check if the object {ObjectKey} exists.
Return the object data or an error.
Make sure that your answer complies with S3 standards, refer to the Amazon S3 documentation for an example.

#### 3. Delete an Object:
HTTP Method: DELETE
Endpoint: /{BucketName}/{ObjectKey}
Behavior:
Verify if the bucket and object exist.
Delete the object and update metadata.
Respond with a 204 No Content status code or an appropriate error message.
Meet the standards.

### Example Scenarios
Scenario 1: Object Upload
A client sends a PUT request to /photos/sunset.png with the binary content of an image.
The server checks if the photos bucket exists, validates the object key sunset.png, and saves the file to data/photos/sunset.png.
The server updates data/photos/objects.csv with metadata for sunset.png and responds with 200 OK.
Scenario 2: Object Retrieval
A client sends a GET request to /photos/sunset.png.
The server checks if the photos bucket exists and if sunset.png exists within the bucket.
The server reads the file from data/photos/sunset.png and returns the binary content with the Content-Type header set to image/png.
Scenario 3: Object Deletion
A client sends a DELETE request to /photos/sunset.png.
The server checks if the photos bucket exists and if sunset.png exists within the bucket.
The server deletes data/photos/sunset.png and removes the corresponding entry from data/photos/objects.csv.
The server responds with 204 No Content.
## Implementation Details:
### Directory Structure:
Use a base directory for storing all data (e.g., data/).
Inside this base directory, create subdirectories for each bucket (data/{bucket-name}/).
Store object files directly in the bucket's directory and maintain a metadata CSV file (objects.csv) to keep track of all objects.
Object Upload Flow:
Bucket Verification: When a PUT request is received, the server checks if the specified bucket exists.
Object Key Validation: The server validates the object key for acceptable characters and length.
Save File: The server writes the binary content to the file system (data/{bucket-name}/{object-key}).
Update Metadata: Update the objects.csv file for the bucket, appending a new entry or updating an existing one.
Error Handling: Handle errors such as insufficient storage, permission issues, and invalid object keys.
Object Retrieval Flow:
Bucket and Object Verification: The server checks if both the bucket and object exist.
Read File: If the object exists, the server reads the file content from disk.
Set Response Headers: The server sets the appropriate MIME type and other headers.
Send Response: The server sends the binary content of the object to the client.
Object Deletion Flow:
Bucket and Object Verification: The server checks if both the bucket and object exist.
Delete File: If the object exists, the server deletes the file from disk.
Update Metadata: Remove the corresponding entry from objects.csv.
Error Handling: Handle cases where the object does not exist or deletion fails due to file system errors.
### Error Handling:
Gracefully handle file access errors (e.g., file not found, permission denied).
Respond with appropriate HTTP status codes for different errors (e.g., 404 Not Found for a missing bucket, 409 Conflict for duplicate bucket names).
## Storing Metadata in a CSV File:
### Bucket CSV File Structure:
Each line in the CSV file represents a bucket's metadata.
The columns could include:
Name: The unique name of the bucket.
CreationTime: The timestamp when the bucket was created.
LastModifiedTime: The last time any modification was made to the bucket.
Status: Indicates whether the bucket is active or marked for deletion.
### Storing Object Metadata in a CSV File:
CSV File Structure for Object Metadata:
Each bucket will have its own object metadata CSV file (e.g., data/{bucket-name}/objects.csv).
The columns could include:
ObjectKey: The unique key or identifier of the object within the bucket.
Size: The size of the object in bytes.
ContentType: The MIME type of the object (e.g., image/png, application/pdf).
LastModified: The timestamp when the object was last uploaded or modified.
## Usage
Your program must be able to print usage information.

Outcomes:

Program prints usage text.
$ ./triple-s --help  
Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory

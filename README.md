## API Service to upload & fetch images

Often we come across scenarios and requirements where we will have to upload and fetch images in our APIs. This repository contains a sample code to achieve the same. The codebase uses Golang & Gin framework to run the server that exposes

- `/upload`
- `/images/:id`

endpoints to upload and fetch the images respectively from the PostgreSQL Database.

### Why only the metadata is stored in the DB?
Glad this interests you! As you can see from the code, the image metadata is stored in the database but the actual image is stored in a folder called `upload`. I did a little bit of digging about this and this is what I found. 

Advantages of storing the metadata in a DB & image in a folder (filesystem):

**Scalability and Performance**: Filesystems are optimized for reading and writing files, making them efficient for storing and serving binary data such as images, videos, or documents. This separation allows you to scale your application by adding more storage as needed without significantly impacting database performance.

**Flexibility**: Storing metadata in a database provides structure and flexibility for searching, querying, and managing your data. You can easily filter, sort, and search for files based on various criteria, such as file name, date, or user who uploaded it.

**Security and Access Control**: Databases provide a robust mechanism for access control and permissions. You can implement fine-grained access control to determine who can view or download specific files. This level of control is often more challenging to achieve with files stored directly in a filesystem.

**Reduced Database Bloat**: Storing large binary files in the database can lead to database bloat and slower performance. By keeping the binary files in the filesystem, you avoid these issues and can keep the database focused on managing metadata.

**Compatibility**: Storing files in the filesystem makes them accessible to various tools and applications outside your application's context, making it easier to interact with the files in different ways.

At the same time there are disadvantages too as below

**Complexity of Synchronization**: Ensuring consistency between metadata in the database and actual files in the filesystem can be challenging.

**Backup and Recovery Complexity**: Managing backups and recovery procedures can be more complex when dealing with both a database and filesystem.

**Data Duplication**: Storing multiple copies of the same file for different users or purposes can lead to data duplication and increased storage requirements.

**Atomicity Across Multiple Operations**: Achieving atomicity across operations involving both metadata and binary files can be challenging.

**Scalability Concerns**: Managing a large number of files in a filesystem as your application scales can become complex.

**Security Risks**: Storing binary files in a filesystem can expose them to security risks without proper access controls.

**Limited Search Capabilities**: Searching for specific content within binary files is more challenging compared to searching metadata.

**Content Management**: Managing a large number of binary files in a filesystem can become unwieldy over time.

**Versioning Complexity**: Maintaining versioning for binary files can be more complex than for metadata.

**Cross-Platform Compatibility**: The choice of filesystem can impact cross-platform compatibility, as different filesystems have varying limitations.

But yea, life is full of options and every option has its advantage & disadvantage. What is important & sensible is to identify the trade offs and to reap the benefits of pros while ignoring with caution the challenges of cons. 

## Usage

1. Clone the repository & navigate inside the folder
    ```bash
    git clone https://github.com/AxiomSamarth/apis-for-img-ops
    cd apis-for-img-ops
    ```

2. Run tidy & vendoring as needed
    ```bash
    go mod tidy
    go mod vendor
    ```

3. Start the server. This starts the server at `http://localhost:8080`
    ```bash
    go run main.go
    ```

4. Either use Postman or the below curl commands to perform the POST & GET operations of images.

    ```bash
    # to upload an image located in a particular path
    curl --location 'http://localhost:8080/upload' \
    --form 'image=@"<path to the image>"'

    # the response will have the reference of the image Id for GET operations
    {
        "Id": 5,
        "message": "Image uploaded successfully"
    }

    # to fetch an image with a particular id
    curl --location 'http://localhost:8080/images/4'

    # the response will be the image
    ```
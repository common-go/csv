# io
- Utilities to load file, save file, zip file
## ETL (Extract-Transform-Load)
Extract-Transform-Load (ETL) is a data integration process involving the extraction of data from various sources, transformation into a suitable format, and loading into a target database or data warehouse.
- Extracting data from various sources.
- Transforming the data into a suitable format/structure.
- Loading the transformed data into a target database or data warehouse.

## Batch processing
- [core-go/io](https://github.com/core-go/io) is designed for batch processing, enabling the development of complex batch applications. It supports operations such as reading, processing, and writing large volumes of data.
- [core-go/io](https://github.com/core-go/io) is not an ETL tool. It provides the necessary libraries for implementing ETL processes. It allows developers to create jobs that extract data from sources, transform it, and load it into destinations, effectively supporting ETL operations.

### Use Cases of [core-go/io](https://github.com/core-go/io) in ETL:
- <b>Data Migration</b>: Moving and transforming data from legacy systems to new systems.
- <b>Data Processing</b>: Handling large-scale data processing tasks like data cleansing and transformation
- <b>Data Warehousing</b>: Loading and transforming data into data warehouses.
### Samples
#### Export from database to file

  ![Export from database to file](https://cdn-images-1.medium.com/max/800/1*IEMXhQXJ0hWZBPL8q2jMNw.png)

##### Samples:
- [go-sql-export](https://github.com/project-samples/go-sql-export): export data from sql to fix-length or csv file.
- [go-hive-export](https://github.com/project-samples/go-hive-export): export data from hive to fix-length or csv file.
- [go-cassandra-export](https://github.com/project-samples/go-cassandra-export): export data from cassandra to fix-length or csv file.
- [go-mongo-export](https://github.com/project-samples/go-mongo-export): export data from mongo to fix-length or csv file.
- [go-firestore-export](https://github.com/project-samples/go-firestore-export): export data from firestore to fix-length or csv file.

#### Import from file to database

  ![Import from file to database](https://cdn-images-1.medium.com/max/800/1*rYaIdKGSd0HwZqZW7pMEiQ.png)
 
  - Detailed flow to import from file to database

    ![Import flow with data validation](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

##### Samples:
- [go-sql-import](https://github.com/project-samples/go-sql-import): import data from fix-length or csv file to sql.
- [go-hive-import](https://github.com/project-samples/go-hive-import): import data from fix-length or csv file to sql.
- [go-cassandra-export](https://github.com/project-samples/go-cassandra-import): import data from fix-length or csv file to sql.
- [go-elasticsearch-import](https://github.com/project-samples/go-elasticsearch-import): import data from fix-length or csv file to sql.
- [go-mongo-export](https://github.com/project-samples/go-mongo-import): import data from fix-length or csv file to sql.
- [go-firestore-export](https://github.com/project-samples/go-firestore-import): import data from fix-length or csv file to sql.

##### Layer Architecture
- Popular for web development

![Layer Architecture](https://cdn-images-1.medium.com/max/800/1*JDYTlK00yg0IlUjZ9-sp7Q.png)

##### Hexagonal Architecture
- Suitable for Import Flow

![Hexagonal Architecture](https://cdn-images-1.medium.com/max/800/1*nMu5_jZJ1omzIB5VK5Lh-w.png)

##### Based on the flow, there are 4 main components (4 main ports):
- Reader, Validator, Transformer, Writer
##### Reader
Reader Adapter Sample: File Reader. We provide 2 file reader adapters:
- Delimiter (CSV format) File Reader
- Fix Length File Reader
##### Validator
- Validator Adapter Sample: Schema Validator
- We provide the Schema validator based on GOLANG Tags
##### Transformer
We provide 2 transformer adapters
- Delimiter Transformer (CSV)
- Fix Length Transformer
##### Writer
We provide many writer adapters:
- SQL Writer: to insert or update data
- SQL Inserter: to insert data
- SQL Updater: to update data

- SQL Stream Writer: to insert or update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- SQL Inserter: to insert data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush. Especially, we build 1 single SQL statement to improve the performance.
- SQL Updater: to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.

- Mongo Writer: to insert or update data
- Mongo Inserter: to insert data
- Mongo Updater: to update data

- Mongo Stream Writer: to insert or update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- Mongo Inserter: to insert data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- Mongo Updater: to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.


## Summary
### File Reader
- File Stream Reader
- Delimiter (CSV format) File Reader
- Fix Length File Reader
### File Writer
- File Stream Writer
#### Delimiter (CSV format) Transformer
- Transform an object to Delimiter (CSV) format
- Transform an object to Fix Length format

## Installation
Please make sure to initialize a Go module before installing core-go/io:

```shell
go get -u github.com/core-go/io
```

Import:
```go
import "github.com/core-go/io"
```

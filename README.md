# Notes

## Architecture
* The code structure is relatively straightforward and in line with most popular go projects.
    * the *internal* directory hosts packages that are meant to be unexposed to packages outside it's ancestral path.
    * the *cmd* directory hosts code responsible for initializing the application and stitching together the underlying abstractions.
* External services and the database are exposed through a layer of abstraction. This serves as a useful logical encapsulation and eases change if the need arises. 
* While not exhaustive the tests included are the ones I needed to get quick feedback on my code.

## Things I should have done
### Data
* The data model is extremely simplistic. A far more robust solution for geo-spacial data would take advantage of a geo-spatial database (or supporting features). The PostGIS extension and it's geometry types allow a greater expressivity while providing a relatively stronger guarantee of correctness. (If I were to implement the additional 3rd task of performing spatial calculations then this would certainly be the way to go.) In the current implementation I store lat/long values as strings so as not to deal with the messiness of floats.
* As a one time project I simply used an SQL script to create the required table. A more future-proof approach would possibly include a migration tool. The fields and tables are also versionless in the current implementaion.
* Bulk/Batch Insert is almost a necessity for this use case and although somewhat trivial to implement now, my initial plan was to version each row (only updating if data values had changed) and so I ended up just leaving the Insert as it is. Both batched inserts and diffed updates would fit this use case well.
* The application should ideally use a tailored role to access the database but in the current implementation simply uses the super user. The constraints on the free tier db combined with the lack of valuable data act as protection. (Ya, I know, but it's really late here!)

### Service
* I had the choice to use the application memory as a cache but chose to keep the service stateless.
* The web server is unauthenticated and the only line of defense is an extremely restrictive global rate limiter (it essentially prevents concurrent requests).
* Polling interval, map size, request timeouts and a host of other values are good candidates for configuration. In the current implementation I use environment variables for addressing external resources.
* The norm is to expose a CLI to the application to help with scripting and automation. For brevity I chose against it.
* I copied and pasted a simple logger I use for my personal applications. Not as powerful as some of the popular open source packages out there but it was fast to setup.
* No robust error handling on the web server. For now the application will just return a HTTP 500 when it runs into problems.

### Map
* I chose to expose a web server to ease reviewer validation. The root path is set to return an image tag with the most recent map.


## Task 3 (Or how I would've done it)
* Although task 3 isn't much of an extension in terms of features, productionizing even a few lines of code is time consuming.
* My first instinct would be to use one of the available packages/extensions for the *Haversine* implementation.
* After qualifying the vehicles based on the criteria, it would be possible to use the *Samsara* API (https://developers.samsara.com/reference#getvehiclesdriverassignments) to gather driver names. If I'm not wrong drivers are ephemeral relative to a vehicle so caching this data may not be plausible (unless a specific schedule was known; easing invalidation). Cache if that is not the case.
* Expose an endpoint through our web server to serve the generated CSV file.


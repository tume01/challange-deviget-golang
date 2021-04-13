# Golang-Challenge

### Solution
The implementation for this challenge was made mainly using go routines. There were two missing
features for this work:
1. Get price by code from cache validating time for eviction
2. Fetch prices from service in a concurrent way

### Get price by code from cache validating time for eviction
For this part, since everything was in memory a good solution was to just delete it from memory. In order
to do it, I added a timer function that will receive the max duration and will clean the map when it ticks.
Under the hood this uses a go routine that will tick when the duration completes. With this approach
we just need to read and write from the cache map and the ticker will handle the eviction.

### Fetch prices from service in a concurrent way
We can fetch prices using go routines since every price doesnt depend on each other. For this 
I used a wait group that will handle go routines sync and will wait on everything to finish before ending the func.
Also, added a mute, so we can be thread safe for reading and writing to the response slice

### Considerations
* With this approach we are assuming memory is not a problem
* The errors returned from getting the price for the item are not reported as the first it comes, this depends on how 
we want to handle this. If we need to get the first error that happens we need to sync with a channel, so we can stop
  other go routines and return the first error encountered.
* A mutex was added to the tests, so we can test using the count param in test, so we know our tests can handle different thread numbers
amounts
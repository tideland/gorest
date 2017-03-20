# Tideland Go REST Server Library

## 2017-03-XX

- Rename internal *envelope* to public *Feedback* in *rest*
- Added *ReadFeedback()* to *Response* in *request*
- Asserts in *restaudit* now internally increase the callstack 
  offset so that the correct test line number is shown
- Added *Response.AssertBodyGrep()* to *restaudit*

## 2017-02-12

- Some renamings in *Request*  and *Response*, sadly
  incompatible to the previous minor release
- More convenience helpers for testing
- Adopted new testing to more packages
- Using http package constants instead of own
  plain strings
- Added documentation to *restaudit*

## 2017-02-10

- Extended *Request* and *Response* of *restaudit* with some
  convenience methods for easier testing
- Adopted *restaudit* changes in *rest* tests

## 2017-01-19

- Renamed type *Query* to *Values*
- Added *Form()* to *Job*

## 2016-12-15

- Added *StatusCode* to feedback envelope
- *JWTAuthorizationHandler* now provides different status codes
  depending on valid tokens, expiration time, and authorization

## 2016-12-09

- *FileServeHandler* now logs the absolute filename and logs
  error if the name is invalid

## 2016-12-06

- *PositiveFeedback()* and *NegativeFeedback()* now also return
  false to be directly used as final return in handler methods

## 2016-12-02

- Added logging to negative responses

## 2016-12-01

- Added missing status code

## 2016-11-23

- Added *JWTFromContext()* to *handlers*
- Later removed JWT context from *handler*; now *jwt* package
  has *NewContext()* and *FromContext()* as usual

## 2016-11-07

- Added *RegisteredHandlers()* to *Multiplexer* retrieve the list
  of registered handlers for one domain and resource
- *Deregister()* is now more flexible in deristering multiple
  or all handlers for one domain and resource at once

## 2016-11-03

- Added *request* package for more convenient requests to REST APIs

## 2016-10-25

- Fixed missing feedback after JWT authorization denial

## 2016-10-24

- Fixed marshalling bug of positive or negative feedback

## 2016-10-18

- Added *Query* type and method for more concenient access to
  query values

## 2016-10-08

- *Job* allows now to enhance its context for following handlers
- *JWTAuthorizationHandler* stores a successfully checked token
  in the job context

## 2016-10-05

- *Formatter.Write()* now also writes the status code

## 2016-10-04

- Improved passing external contexts into an environment, e.g.
  containing database connection pools
- Changed multiplexer configuration to now use *etc.Etc* from
  the *Tideland Go Library*
- More robust basepath handling now

## 2016-09-29

- Fixed bug with public handler types

## 2016-09-27

- Added methods for the lazy loading and rendering of templates
- Sadly has little impact on the rendering interface

## 2016-09-19

- Finished rework after adding of JSON Web Token package

## 2016-08-21

- Migrated *Tideland Go Library* web package after some rework
  into this new project

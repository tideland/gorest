# Tideland Go REST Server Library

## 2016-11-23

- Added *JWTFromContext()* to *handlers*

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

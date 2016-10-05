# Tideland Go REST Server Library

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

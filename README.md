# FarmEasy

## Problem statement

FarmEasy is an online portal for farmers to lend and borrow farming instruments. Consider a scenario where a farmer has sugarcane planted on his field and at the time of harvesting, farmers would need a sugarcane harvester machine, but the thing is that the machine costs more than 25 lakhs in the market and a small farmer doesn't have that kind of money.

On the other hand, a farmer who owns that sugarcane harvesting machine, might have that machine just sitting around and he might wanna earn some money by lending that machine to some other farmer at an hourly charge.

FarmEasy provides the solution to this problem.

## Roles

#### Farmer :

A user is farmer who can list his farming instruments on the portal to rent and also it can be a farmer who wants to borrow the instrument

## Features

- farmer is able to list his machines on the portal, along with the base hourly charge
- farmer is able to see all the machines listed on the portal
- farmer is able to rent the machine from the web application
- farmer is able to get the invoice for booking

## [Api specification](https://docs.google.com/document/d/1LWpB_4gvnwUaYkubR511bj_4SYm0HcvGXmhwwKXOBBQ/edit?usp=sharing)

## DB schema

![db-schema](/db-schema.png)

## How to start

- Download and setup Postgresql

- create farmeasy database

- Setup DB URI in application.yml file

```console
DB_URI: "postgresql://username:password@localhost:5432/farmeasy?sslmode=disable"
```

- execute

```console
go build
./FarmEasy migrate
./FarmEasy start
```

- To run testcases

```console
go test ./...
```

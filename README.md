# Guild Project DDD

## Project Overview

The objetive of this project is design a kind of guild system, popular in rpg games, using Go and trying to put DDD practices as long the project evolves.
Trying to push my understand of DDD, I'm focus to get maximum modularity, where each module will acting as a separated microservices.
I believe this approach will help me face some problems, like authentication, consistency and other design problems to think and find way to address it.

## High Level Architecture

![image info](./pictures/guild-overview.png)

### Key Modules

- Auth Server - Responsible to handle authentication and user session. In the future I want to connect it to a Api Gateway to centralize the authentication/authorization concerns. I also want to change this to a prebuilt tool like Keycloack and see how these will impact my system design.

- Player = Responsible to handle all players concerns. Here we want to be able to create new players. Here I intent to explore data synchronization problem where, wherenever a player is created, it should have a Login in Auth Server. I want to explore a sync approach, and an async approach using domain events and see the challenge to archived the small latency possible to get eventual consistency.

- Guild - Responsible to handle all guild concerns. Here we want to be able to create guild, invite others player to our guild, promote, demote or kick players. Players will be able to donate gold and cash to a Guild and we want to track those actions to create a report of transactions.

- Vault - Responsible to handle all vault concerns. Although a Vault have its own life cycle, to be consistent it should communicate with Player and Guild module. Like in the Player-Login problem, I want to explore both sync and async approachs, using domain events and push further my understand about Domain Services.

### Key Features to explore in future

- Api Gateway
- Service Discovery
- Rate Limit
- Service Mesh
- Circuit Breaker

## Progress

### Auth Server

**Status:** _in progress_

**Features Available:**

- Create a Login
- Login
- Renew Token and Session
- Revoke session

**To Do List:**

- [x] Create basic structure
- [ ] Change database mocks to "faker"approach and improve unit tests consistency
- [ ] Improve renew token business logic
- [ ] Create an Auth Middleware to protect /revoke and /renew endpoints
- [ ] Write tests to app services and route
- [ ] Write tests to infra layer components

# Design Specification

# Purpose and scope
An app which allows users to share and discover quotes. Quotes can be from any source, be it book, article, video, or spoken word. The goal is for everyone to share the ideas that inspire them as well as get inspired by others' ideas.

Users should be able to register an account, activate it and post quotes. Users can navigate to the discover page, which displays all the quotes posted by activated users. Quotes can be searched in a search bar. 

# Requirements
### Frontend
- [ ] "My quotes" page when user is logged in
- [ ] Discover page which gives a list of 10 quotes by time posted or most liked
- [ ] Pagination on the discover page
- [ ] Search bar on discover page
- [ ] UI is responsive to mobile screens
### [[Quotable backend|Backend]]

- [ ] User can register an account with username and email
- [ ] User will be sent an activation email
- [ ] Anonymous users can view the discover page but cannot create quotes
- [ ] Users that have not been activated can create quotes but those quotes cannot be viewed in the discover page
- [ ] Users should be able to view their quotes in the quotes page
- [ ] Users can tag quotes with categories to make them easier to search for later.
- [ ] Users can like posts
- [ ] Users can sort posts
# Data model
We will be using a relational database.

Users table:
- user_id PK
- username
- email
- hashed_password
quotes table:
- quotes_id PK
- content
- author
- tags
likes table:
- user_id FK
- quote_id FK
# Routes
A plan for the set of routes is shown below

| Method | URL                    | Action                                   |
| ------ | ---------------------- | ---------------------------------------- |
| POST   | /user/register         | Register the user                        |
| GET    | /quotes                | Query the quotes using url query params  |
| GET    | /quotes/:id            | Query quote by quote ID                  |
| POST   | /tokens/auth           | Create an auth token for the user        |
| GET    | /users/:user_id/quotes | Query the quotes of user with id user_id |
| POST   | /users/password        | Change password of user                  |
| PATCH  | /quotes/:id            | Update the quote                         |
| DELETE | /quotes/:id            | Delete the quote                         |

# Tech stack
### Prototype
For the prototype I will be using Go as the backend language and Postgres as the database. The pages will be server-side rendered using Go templates. I plan to use minimal JS and CSS in regards to the front end.
### Future plans
Use of react or svelte to create a proper front-end UI.
### Goal of project:

This api is built to support the React Native project [here](https://github.com/willybeans/freedu_rn). It is designed to provide users the ability to create and share language learning content with one another. Most importantly it provides a space where someone can create their own language learning content themselves. Stephen Krashens theory of language learning is the driving influence behind the design and features that have been implemented in this application

#### What problem is being solved?

- As a language learner, it is difficult to find viable content to practice reading
- Leveraging real world content in particular is often difficult because consistently translating content manually is time consuming
- It is hard to keep track of what words one is currently learning, or have learned, which robs the learner of a concrete sense of progress
- Sharing content with other people, and also using other peoples content, is often a more efficient than always having to look for content yourself. Fluent speakers can often cater their speech towards your language proficiency level easier than it is to find books that 'fit' with your level.
- Creating a social aspect behind the often solitary process of language learning helps build relationships with other language learners which foster motivation, and often dispells fears
- When making the transition from learning to speaking, it is often hard to find content/tools that eases you into that transition.

### What Stack? What is where?

- Golang API REST API with websockets
- using Chi for routing
- gorilla/websockets for websocket connection
- home grown PUB / SUB logic for websocket subscriptions
- Leverages tesseract wrapper `gosseract` for page scraping
- Using `PostgreSQL` for DB
- Uses logging middleware from `zerolog`

## Contribution Guidlines:

## Set Up:

1. Fork on github for cleanliness
2. Clone your fork locally (make sure to [pull from this parent project](https://stackoverflow.com/questions/13828230/pulling-changes-from-fork-parent-in-git) instead of your branch, otherwise you wont get the current git state)

   - `git clone git@github.com:<my-github-name>/freedu_go.git`

3. [Setup Golang](https://go.dev/doc/install) and also install the vscode go language support extension
4. setup PostgreSQL with either [downloading](https://www.postgresql.org/download/) or homebrew
5. You will need to add a `.env` file to the root of your project in order to allow the server to connect with psql. you can copy the commented out values from `./create_database.sh` into a `.env` file or you can write your own. Making sure to replace "your_database_name", "your_username", and "your_password" with your actual PostgreSQL database credentials, ensure that you have the psql command-line tool installed on your system
6. We now have to configure the `create_database.sh` file, and make it executable by running `chmod +x create_database.sh` in the terminal. Then you can execute the script from the root using `./create_database.sh` which will populate your psql instance with the appropriate database and tables
7. If you intend on using the tesseract endpoint, you will need to follow the [setup instructions](https://github.com/otiai10/gosseract?tab=readme-ov-file#installation) for a local version of tesseract, otherwise the page scraping will not work. If you want to skip this step, you may simply comment out the handler/route logic using goserract
8. Install air to allow for hotreloading by pasting this into your terminal `go install github.com/cosmtrek/air@latest
`

## To Run

- making sure psql is running, open a terminal in root dir, and run the command `LOG_LEVEL=1 APP_ENV="development" air` for prettified logging
- `LOG_LEVEL=1 APP_ENV="production" air` for production
  -If you want to see the DEBUG and TRACE logs, quit the air command with Ctrl-C and rerun it with LOG_LEVEL set to -1: `LOG_LEVEL=-1 APP_ENV=development air`

### Your first pull request

- Make clear your intention to work on a problem in the issue section by either:

1.  Making an issue yourself and leaving a comment of your intent to complete the issue
2.  Comment on existing issue with your intention to fix it

- All code should be compliant to the proper lint rules
- Use a branch naming convention like `fix/short-fix-description` or `feature/short-feature-description`
- Please keep all push requests concise
- _Avoid_ pushing more than one file at a time (avoid `git add .` unless you are certain it is not adding additional unrelated material)
- Always be up to date: `git pull` to avoid a merge conflict avalanche before your PR
- Push to staging branch so your commit can be tested and confirmed

#### Issues

- Pick an unassigned issue that you think you can accomplish, add a comment that you are attempting to do it.
- Feel free to propose issues that aren’t described! Get the okay that it is inline with the project goals before working.

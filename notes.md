Basic version:

Someone creates a poll, then sends link to friends
- Specifies food options (start w/ plain text, add something like autocompelte later)

Friends click link, rank food choices from 1-x, click submit

Server will process results as they come in, declare a winner

The creator of the poll will have a way to see results, and once they see the correct number of votes, they can show the winner

TYPES

Poll
- ID
- Options
- Creator slug (is the link)
- Voter slug (is the link)

Results
- ID
- Poll ID
- results

CURRENT STATUS:
- Polls get inserted into sqlite database
- Next, make views for Creator's view (non-functional), and Poll View
- Make URLs w/ slugs link to appropriate views

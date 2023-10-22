# flashwise backend

Backend for the Flashwise app.

## Testing

1. Create a user

   ```
   http POST localhost:8080/register \
       username="charlie" \
       email="charlie.root@gmail.com" \
       password1="tiger123" \
       password2="tiger123"
   ```

1. Create a flashcard set:

   ```
   http POST localhost:8080/sets \
       'Authorization: Token 6d0c1a5ecb334e176c5d13e8d24c282a8b45684d' \
       status:=true \
       description='Geography'
   ```

1. Get user's flashcard sets:

   ```
   http GET localhost:8080/sets \
       'Authorization: Token 6d0c1a5ecb334e176c5d13e8d24c282a8b45684d'
   ```

1. Create a flashcard

   ```
   http POST localhost:8080/flashcards \
       'Authorization: Token 6d0c1a5ecb334e176c5d13e8d24c282a8b45684d' \
       front="Capital of Poland" \
       back="Warsaw" \
       flashcard_set:=1
   ```

1. Get user's flashcards:

   ```
   http GET localhost:8080/flashcards \
       'Authorization: Token 6d0c1a5ecb334e176c5d13e8d24c282a8b45684d'
   ```

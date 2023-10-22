# flashwise backend

Backend for the Flashwise app.

## Testing

1. Create a user

1. Create a flashcard set:

   ```
   http POST localhost:8080/sets \
       'Authorization: Token c8dd1617cf1f35965cb8d2f827f4c2d834f2958b' \
       status:=true \
       description='Geography'
   ```

1. Get user's flashcard sets:

   ```
   http GET localhost:8080/sets \
       'Authorization: Token c8dd1617cf1f35965cb8d2f827f4c2d834f2958b'
   ```

1. Create a flashcard

   ```
   http POST localhost:8080/flashcards \
       'Authorization: Token c8dd1617cf1f35965cb8d2f827f4c2d834f2958b' \
       front="Capital of Poland" \
       back="Warsaw" \
       set_id:=1 \
       author_id:=1
   ```

# flashwise backend

Backend for the Flashwise app.

## Testing

### Create a flashcard

```
http POST localhost:8080/flashcards \
    front="Front2" \
    back="Back2" \
    set_id:=1 \
    author_id:=2
```

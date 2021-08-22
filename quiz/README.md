The official solution can be found at https://github.com/gophercises/quiz/tree/solution-p2

## Insights
- I have used `panic` instead of `exit`. Actually, `exit` is a better solution since those calls do not indicate programmer errors
- I didn't clarify the higher level cause when an error occurred
- They trimmed the answers when parsing, while I did it when comparing with the input. Since they were used only once, it does not really matter in this case
- They checked for the channels in each loop iteration, while I have iterated the problems within a goroutine and checked for the channels afterwards. Their solution looks more elegant.
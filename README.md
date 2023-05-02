# Bulk SMS sender

The code is used to send SMS in bulk through MessageBird.

To use it:

- Create an account on MessageBird and have your API key ready
- Create a `.csv` file with the data you want to use to send SMSs. The `.csv` file should contain the same number of columns for all rows, should be comma separated, and should **not** have a header. It should contain the recipient phone number (one per row)
- Copy `config.yml.template` to `config.yml`
    - Customize the configuration keys such as `message` and the `text`. 
    - Under the `columns` key, write the columns contained in the `.csv` file, in the right order.
    - Under the `phoneColumns` key, write the column name that contains the phone number
    - For the `message`, you can use data contained in the `.csv` by referencing it with `{{ column_name }}`
- Set your MessageBird API into the `MESSAGEBIRD_API` environment variable and then run the code with `go run main.go`


## Example

- `config.yml`:

```yaml
message: "Hey {{ name }}! Today is your lucky day! Good luck!"
sender: "GitHub"
filename: "names.csv"
columns:
  - name
  - lastname
  - phone
phoneColumn: phone
```

- `names.csv`:

```
John,Doe,+31650000001
Jane,Doe,+31650000002
```

- (Truncated) Output of `go run main.go`

```
2023/05/02 21:16:58 &{644bf https://rest.messagebird.com/messages/644bf mt sms GitHub Hey John! Today is your lucky day! Good luck!  <nil> [{31650000001 sent 2023-05-02 19:16:58}]}}
```

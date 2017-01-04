# t

Japanese <-> English translate in terminal

## Require

### Use Watson Language Translator

- [Service credentials](https://console.ng.bluemix.net/dashboard/apps/)
  - username (different from your Bluemix account)
    - `export T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME=<Your Watson Language Translator API username>`
  - password (different from your Bluemix account)
    - `export T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD=<Your Watson Language Translator API password>`

## Install

You can download binary from [release page](https://github.com/unblee/t/releases) and place it in `$PATH` directory.

## Usage

```
t [input text]

or

echo [input text] | t

  option:
    -v --version: show version

  t translates input text specified by argument or STDIN using Watson Language Translation API.
  Source language will be automatically detected.

  export T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME=<Your Watson Language Translator API username>
  export T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD=<Your Watson Language Translator API password>

  Example:
    $ t Good morning!
    おはようございます!
    $ t おはようございます!
    Good morning!
```

## Note

- No available in `mintty` :sob: (isatty() not working)
- The text must be the **UTF-8** when using a pipe input
  - exec `echo "UTF-8 text" | t`

## Supported

- [IBM Bluemix Watson Language Translator](https://www.ibm.com/watson/developercloud/language-translator.html)

## Todo

- [ ] Add unit tests
- [ ] [Support Google Translate API](https://cloud.google.com/translate/)
- [ ] [Support Microsoft Translator Text API](https://www.microsoft.com/cognitive-services/en-us/translator-api)

## Development

### Require

- unix shell(e.g. bash, zsh)
  - msys2 on windows (or WSL)
- golang >= 1.7
- make
- [glide](https://github.com/Masterminds/glide)
- docker
- editor/IDE with editorconfig plugin
- github token

## Reference

Many thanks!

- [みんなのGo言語](http://gihyo.jp/book/2016/978-4-7741-8392-3)
- [https://github.com/haya14busa/gtrans](https://github.com/haya14busa/gtrans)
- [http://deeeet.com/writing/2016/11/01/go-api-client/](http://deeeet.com/writing/2016/11/01/go-api-client/)
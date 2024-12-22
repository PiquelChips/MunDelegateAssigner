# MunDelegateAssigner

Will randomly assign each delegate a country for an MUN debate. If you keep the `history.csv` file, delegates who have recently had an important country will have a lower chance of getting one again.
Source code can be found on [GitHub](github.com/PiquelChips/MunDelegateAssigner).

## TODO

- Remove need for country to have a value

## How to use

### Setup

- Make sure the configuration file `config.ini` has every country you need. Each category has its own importance which is defined in the general section. The salt value is for the amount of randomness you want.
- Each country must be assigned a value (<country>=0). The value does not matter it just needs to be there
- Make sure to have a `delegates.csv` file which contains the list of delegates that need countries (including chairs). Each line must have one delegate and no spaces are allowed.

### Assigning

- Right-click the folder containing the assigner and select `Open in Terminal`.
- In the terminal, type `.\DelegateAssigner` and press enter.
- Follow the instructions to enter the names of the chairs.
- The assignments will be saved to `assignments.csv`. You can then import this file into `Google Sheets`, `Excel`, or any other spreadsheet software.
- In most software you can import by selecting `file` then `import` or `import .csv` and then selecting the `delegates.csv` file.
- Once imported, you will get two columns, one with the delegates and one with their assigned country

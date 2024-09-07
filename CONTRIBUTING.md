# Contributing to go-backend-template

First off, thanks for taking the time to contribute!

All types of contributions are encouraged and valued.
See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them.
Please make sure to read the relevant section before making your contribution.
It will make it a lot easier for us maintainers and smooth out the experience for all involved.
The community looks forward to your contributions.

> And if you like the project, but just don't have time to contribute, that's fine.
There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues


## Table of Contents
- [Reporting an Issue](#reporting-an-issue)
- [I Have a Question](#i-have-a-question)
- [Contributing](#contributing)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Contributing](#contributing)
- [Convenstions](#conventions)
- [Attributions](#attributions)


## Reporting an Issue

A great way to contribute to the project is to send a detailed report when you encounter an issue.
We always appreciate a well-written, thorough bug report, and will thank you for it!

Check that our [issue](https://github.com/npushpakumara/go-backend-template/issues) database doesn't already include that problem or suggestion before submitting an issue.
If you find a match, you can use the "subscribe" button to get notified on updates.
Do not leave random "+1" or "I have this too" comments, as they only clutter the discussion, and don't help resolving it.
However, if you have ways to reproduce the issue or have additional information that may help resolving the issue, leave a comment.

When reporting issues, always include:

- Stack trace (Traceback)
- OS, Platform and Version (Windows, Linux, macOS, x86, ARM)
- Possibly your input and the output
- Explain the behavior you would expect and the actual behavior.
- Steps required to reproduce the problem if possible and applicable so that someone else can follow to recreate the issue on their own.
This usually includes your code. 
- When sending lengthy log-files, consider posting them as a [gist](https://gist.github.com).
Don't forget to remove sensitive data from your logfiles before posting (you can replace those parts with "REDACTED").


## I Have a Question

> If you want to ask a question, we assume that you have read the available [README](https://github.com/npushpakumara/go-backend-template/blob/main/README.md).

Before you ask a question, it is best to search for existing [Issues](https://github.com/npushpakumara/go-backend-template/issues) that might help you.
In case you have found a suitable issue and still need clarification, you can write your question in this issue.
It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:
- Open an [Issue](https://github.com/npushpakumara/go-backend-template/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (go, OS, etc), depending on what seems relevant.


## Contributing

Not sure if that typo is worth a pull request ? Found a bug and know how to fix it ? Do it ! We will appreciate it.
Any significant improvement should be documented as a GitHub issue before anybody starts working on it.


> ### Legal Notice 
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.


## Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for go-backend-template, **including completely new features and minor improvements to existing functionality**.
Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.


### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [README](https://github.com/npushpakumara/go-backend-template/blob/main/README.md) carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/npushpakumara/go-backend-template/issues) to see if the enhancement has already been suggested.
If it has, add a comment to the existing issue instead of opening a new one.
- Find out whether your idea fits with the scope and aims of the project.
It's up to you to make a strong case to convince the project's developers of the merits of this feature.
Keep in mind that we will prioritize features that will be useful to the majority of our users and not just a small subset.


### Submitting a Good Enhancement Suggestion

Enhancement suggestions are tracked as [GitHub issues](https://github.com/npushpakumara/go-backend-template/issues).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why.
At this point you can also tell which alternatives do not work for you.
- You may want to **include screenshots and animated GIFs** which help you demonstrate the steps or point out the part which the suggestion is related to.
You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
- **Explain why this enhancement would be useful** to most go-backend-template users. 
ou may also want to point out the other projects that solved it better and which could serve as inspiration.


## Conventions

Fork the repository and make changes on your fork in a feature branch:

- If it's a bug fix branch, name it XXXX-something where XXXX is the number of the issue.
- If it's a feature branch, create an enhancement issue to announce your intentions, and name it XXXX-something where XXXX is the number of the issue.

Submit unit tests for your changes.
Go has a great test framework built in; use it!
Update the documentation when creating or modifying features.

Write clean code.
Universally formatted code promotes ease of writing, reading, and maintenance.
Always run `gofmt -s -w file.go` on each changed file before committing your changes.
Most editors have plug-ins that do this automatically.

Commit messages must must be under 50 chars written in the imperative form.
This can be followed by an optional, more detailed explanatory text which is separated from the summary by an empty line.

Code review comments may be added to your pull request.
Discuss, then make the suggested modifications and push additional commits to your feature branch.
New commits show up in the pull request automatically, but the reviewers are notified only when you comment.

Before you make a pull request, squash your commits into logical units of work. 
A logical unit of work is a consistent set of patches that should be reviewed together.\
**Example:** Upgrading the version of a vendored dependency and taking advantage of its now available new feature constitute two separate units of work.
Implementing a new function and calling it in another file constitute a single logical unit of work.
The very high majority of submissions should have a single commit, so if in doubt: squash down to one.

After every commit, make sure the test suite passes.
Include documentation changes in the same pull request so that a revert would remove all traces of the feature or fix.

Include an issue reference like Closes #XXXX or Fixes #XXXX in the pull request description that close an issue.
Including references automatically closes the issue on a merge.


## Attributions

This Code of conduct is inspired by [Docker CLI CONTRIBUTING.md](https://github.com/docker/cli/blob/master/CONTRIBUTING.md) and [contributing.md](https://contributing.md/)

# Note

This CLI tool is useful for generating various markdown templates, prompting the user for required information. Right now it assumes the user has vscode installed and uses that for editing.

## Usage

To use note, just make sure the executable is somewhere on your classpath. Then run `note <command>`. The application will prompt you for details required for the given template, produce the template in your configured directory and then open vscode at the location.

The following commands are available:

**blog** - creates a blog post with the current date and some useful metadata at the top  
**book** - sets up a directory and series of files ready for detailed book notes  
**meet** - useful for quickly putting together notes for a meeting  
**morning** - for creating morning pages  

## Configuration

Create a configuration file either in your home directory called `.note.yml` (with a prepended period), or in your .config directory called `note.yml` (without a prepended period).

```yaml
fullname: Your Full Name
paths:
  base: /home/your/notes
  blog: /home/your/blog
  book: /home/your/notes/books
  morning: /home/your/notes/morning-pages
  meeting: /home/your/notes/meetings
```

## Building the Project

Before building you need to make sure to package the templates:

```shell
pkger -include /templates -o /cmd
```

Then you can build the executable as you normally would:

```shell
go install .
```
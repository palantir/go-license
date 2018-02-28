amalgomate-plugin
=================
`amalgomate-plugin` is a [godel](https://github.com/palantir/godel) plugin for amalgomate. The plugin runs `amalgomate`
based on provided configuration. It also runs as part of the `--verify` task and verifies that the running the task
would not alter the content of the output directory.

Tasks
-----
* `amalgomate`: runs amalgomation. Runs for all of the entries specified in the configuration in order. The working
  directory is set to be the project directory.

Verify
------
When run as part of verification that does not apply, the task fails if running the task would alter any of the contents
of the output directory. The output directory will not be modified as part of this process. Note that the implementation
operates by temporarily making a copy of the output directory. For this reason, one should avoid having large files in
the output directory.

Config
------
The configuration for this plugin is in a file called `amalgomate.yml`. The configuration should be of the following
form:

```yaml
amalgomators:
  okgo:
    config: ./okgo/checks.yml
    output-dir: ./okgo/generated_src
    pkg: amalgomatedformatters
```

The top-level `amalgomators` is a map where the key is the name of the amalgomation task and the value is the
configuration for that task. `config` specifies the location of the `amalgomate` configuration file for the task, while
`output-dir` specifies the location of the directory to which the output should be written. Both of these are specified
as relative paths relative to the project directory. The `pkg` parameter specifies the package name of the generated
`amalgomate` code.

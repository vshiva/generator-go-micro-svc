'use strict';
const Generator = require('yeoman-generator');
const chalk = require('chalk');
const yosay = require('yosay');
const path = require('path');
const mkdir = require('mkdirp');

module.exports = class extends Generator {
  paths() {
    this.sourceRoot(path.join(__dirname, './templates/__service_name__'));
    this.destinationRoot(process.env.GOPATH || './');
  }

  prompting() {
    const cb = this.async();

    // Have Yeoman greet the user.
    this.log(yosay(`${chalk.red('Go Lang Micro Service')}`));

    let prompts = [
      {
        type: 'input',
        name: 'serviceName',
        message: `What is the name of your ${chalk.yellow('service')}?`,
        default: 'mysvc'
      },
      {
        type: 'input',
        name: 'repo',
        message: `What is your URL ${chalk.yellow('repository')}?`,
        default: 'github.com'
      },
      {
        type: 'input',
        name: 'repoUsr',
        message: `What is your ${chalk.yellow('User or Org')} name of the repository?`,
        default: 'me'
      },
      {
        type: 'confirm',
        name: 'vendor',
        message: `Would you like to commit ${chalk.yellow('vendor')}?`,
        default: false,
        store: true
      }
    ];

    return this.prompt(prompts).then(props => {
      // To access props later use this.props.someAnswer;
      this.props = props;
      const svcName = props.serviceName.replace(/\s+/g, '');
      this.templateData = {
        serviceName: svcName.toLowerCase(),
        servicePName: svcName,
        repoUrl: props.repo + '/' + props.repoUsr + '/' + svcName,
        vendor: props.vendor,
        licenseText: ''
      };
      cb();
    });
  }

  async writing() {
    console.log('Generating tree folders');
    console.log(this.templateData);
    let pkgDir = this.destinationPath('pkg');
    let srcDir = this.destinationPath(path.join('src/', this.templateData.repoUrl));
    let binDir = this.destinationPath('bin');

    mkdir.sync(pkgDir);
    mkdir.sync(srcDir);
    mkdir.sync(binDir);

    this.fs.copyTpl(this.sourceRoot() + '/*', srcDir, this.templateData);

    this.fs.copyTpl(this.sourceRoot() + '/.*', srcDir, this.templateData);

    this.fs.copyTpl(
      this.sourceRoot() + '/.vscode/**',
      path.join(srcDir, '.vscode'),
      this.templateData
    );

    this.fs.copyTpl(
      this.sourceRoot() + '/pkg/**/*',
      path.join(srcDir, 'pkg'),
      this.templateData
    );

    this.fs.copyTpl(
      this.sourceRoot() + '/cmd/**/*',
      path.join(srcDir, 'cmd'),
      this.templateData
    );

    this.fs.copyTpl(
      this.sourceRoot() + '/deployment/**',
      path.join(srcDir, 'deployment'),
      this.templateData
    );

    this.fs.copyTpl(
      this.sourceRoot() + '/vendor/**',
      path.join(srcDir, 'vendor'),
      this.templateData
    );

    this.fs.move(
      path.join(srcDir, '__svc_name__.proto'),
      path.join(srcDir, this.templateData.serviceName + '.proto')
    );

    this.fs.move(
      path.join(srcDir, 'pkg', '__svc_name__pb', '*'),
      path.join(srcDir, 'pkg', this.templateData.serviceName + 'pb')
    );

    this.fs.delete(path.join(srcDir, 'pkg', '__svc_name__pb'));
  }
};

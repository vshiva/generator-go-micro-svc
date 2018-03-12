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
        message: 'What is the name of your service ?',
        default: 'mysvc'
      },
      {
        type: 'input',
        name: 'repo',
        message: 'What is your URL repository ?',
        default: 'github.com'
      },
      {
        type: 'input',
        name: 'repoUsr',
        message: 'What is your User or Org name of the repository ?',
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
        licenseText: ''
      };
      cb();
    });
  }

  async writing() {
    console.log('Generating tree folders');
    let pkgDir = this.destinationPath('pkg');
    let srcDir = this.destinationPath(path.join('src/', this.repoUrl));
    let binDir = this.destinationPath('bin');

    mkdir.sync(pkgDir);
    mkdir.sync(srcDir);
    mkdir.sync(binDir);

    this.fs.copyTpl(this.sourceRoot() + '/*', '/', this.templateData);
    this.fs.copyTpl(
      this.sourceRoot() + '/__svc_name__pb/**',
      '/' + this.templateData.serviceName + '/',
      this.templateData
    );
    this.fs.copyTpl(this.sourceRoot() + '/.vscode/**', '/.vscode/', this.templateData);
    this.fs.copyTpl(
      this.sourceRoot() + '/deployment/**',
      '/deployment/',
      this.templateData
    );
    this.fs.copyTpl(this.sourceRoot() + '/server/**', '/server/', this.templateData);
    this.fs.copyTpl(this.sourceRoot() + '/state/**', '/state/', this.templateData);
    this.fs.copyTpl(this.sourceRoot() + '/vendor/**', '/vendor/', this.templateData);
    this.fs.move(
      this.sourceRoot() + '/__svc_name__.proto',
      this.templateData.serviceName + '.proto'
    );
  }
};

'use strict';
const Generator = require('yeoman-generator');
const chalk = require('chalk');
const yosay = require('yosay');
const path = require('path');
const mkdir = require('mkdirp');
const camelCase = require('camelcase');

const prompts = [
  {
    name: 'serviceName',
    message: `What is the name of your ${chalk.red('service')}?`,
    default: 'mysvc'
  },
  {
    type: 'list',
    name: 'repo',
    message: `What is your ${chalk.red('source repository')}?`,
    choices: [
      {
        name: 'github.com',
        value: 'github.com'
      },
      {
        name: 'bitbucket.com',
        value: 'bitbucket.com'
      }
    ],
    default: 'github.com',
    store: true
  },
  {
    type: 'input',
    name: 'repoUsr',
    message: `What is your ${chalk.red('User or Org')} name of the source repository?`,
    default: process.env.USER || 'me',
    store: true
  },
  {
    type: 'confirm',
    name: 'vendor',
    message: `Would you like to commit ${chalk.red('vendor')}?`,
    default: false,
    store: true
  },
  {
    type: 'list',
    name: 'containerRegistry',
    message: `What is your ${chalk.red('container registry')}?`,
    choices: [
      {
        name: 'ocir.io',
        value: 'ocir.io'
      },
      {
        name: 'docker.io',
        value: 'docker.io'
      },
      {
        name: 'quay.io',
        value: 'quay.io'
      },
      {
        name: 'gcr.io',
        value: 'gcr.io'
      }
    ],
    default: 'ocir.io',
    store: true
  },
  {
    type: 'input',
    name: 'repoName',
    message: `What is the ${chalk.red('User or Org')} name of the ${chalk.red(
      'container image repository'
    )} ?`,
    default: process.env.USER || 'me',
    store: true
  }
];

const sanitizeServiceName = serviceName => {
  return serviceName
    .trim()
    .replace(/  +/g, ' ')
    .split(' ')
    .join('-')
    .replace('_', '-')
    .replace(/[^0-9a-z-]/gi, '')
    .replace(/-+/g, '-')
    .toLowerCase();
};

module.exports = class extends Generator {
  paths() {
    this.sourceRoot(path.join(__dirname, './templates/__service_name__'));
    this.destinationRoot(process.env.GOPATH || './');
  }

  prompting() {
    const cb = this.async();

    // Have Yeoman greet the user.
    this.log(yosay(`${chalk.red('Go Lang Micro Service')}`));
    return this.prompt(prompts).then(props => {
      // To access props later use this.props.someAnswer;
      this.props = props;
      const pkgName = sanitizeServiceName(props.serviceName);
      const svcName = camelCase(pkgName);
      this.templateData = {
        serviceName: svcName,
        serviceLName: svcName.toLowerCase(),
        servicePName: svcName.charAt(0).toUpperCase() + svcName.slice(1),
        repoUrl: props.repo + '/' + props.repoUsr + '/' + pkgName,
        vendor: props.vendor,
        pkgName: pkgName,
        containerRegistry: props.containerRegistry,
        repoName: sanitizeServiceName(props.repoName),
        licenseText: ''
      };
      cb();
    });
  }

  async writing() {
    console.log('Generating tree folders');
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
      this.sourceRoot() + '/helm/**',
      path.join(srcDir, 'helm'),
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
      path.join(srcDir, 'cmd', '__service_name__', '*'),
      path.join(srcDir, 'cmd', this.templateData.serviceName)
    );

    this.fs.delete(path.join(srcDir, 'cmd', '__service_name__'));
  }
};

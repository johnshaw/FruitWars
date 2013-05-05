module.exports = (grunt) ->
  
  grunt.initConfig
    pkg: grunt.file.readJSON('package.json')

    minispade:
      options:
        renameRequire: true
        useStrict: false
        prefixToRemove: 'js/'
      files:
        src: ['js/**/*.js']
        dest: 'fruits.js'
        
    regarde:
      js:
        files: 'js/**/*.js'
        tasks: ['livereload', 'regarde']

  grunt.loadNpmTasks('grunt-contrib-livereload')
  grunt.loadNpmTasks('grunt-regarde')

  grunt.registerTask('default', ['livereload-start', 'regarde'])

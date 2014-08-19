var gulp = require('gulp');
var sass = require('gulp-sass');
var csso = require('gulp-csso');
var uglify = require('gulp-uglify');
var concat = require('gulp-concat');
var plumber = require('gulp-plumber');
var coffee = require('gulp-coffee');
var gutil = require('gulp-util');
var clean = require('gulp-clean');

var paths = {
  fonts: 'bower_components/font-awesome/fonts/*'
};

var cssFiles = [
    'bower_components/rickshaw/rickshaw.min.css'
];

var filesToClean = [
    'public/javascript/*.js',
    'public/stylesheets/*.css',
    'public/fonts'
];

gulp.task('clean', function(cb) {
    gulp.src(filesToClean, {read: false})
        .pipe(clean());
});

gulp.task('sass', function() {
  gulp.src('public/stylesheets/styles.scss')
    .pipe(plumber())
    .pipe(sass())
    .pipe(csso())
    .pipe(gulp.dest('public/stylesheets'));
});

gulp.task('compress', function() {
  gulp.src([
    'bower_components/jquery/dist/jquery.min.js',
    'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap.js',
    'bower_components/rickshaw/vendor/d3.min.js',
    'bower_components/rickshaw/vendor/d3.layout.min.js',
    'bower_components/rickshaw/rickshaw.js',
    'bower_components/lodash/dist/lodash.min.js',
    'public/javascript/*.js',
    '!public/javascript/app.min.js'
  ])
    .pipe(concat('app.min.js'))
    .pipe(gulp.dest('public/javascript'));
});

gulp.task('coffee', function() {
  gulp.src('./public/javascript/*.coffee')
    .pipe(coffee({bare: true}).on('error', gutil.log))
    .pipe(gulp.dest('./public/javascript'));
});

// Copy all fonts to the public folder
gulp.task('fonts', function() {
  gulp.src(paths.fonts)
    .pipe(gulp.dest('public/fonts'));
});

// Copy all additional required CSS files to the public folder
gulp.task('css', function() {
  gulp.src(cssFiles)
    .pipe(gulp.dest('public/stylesheets'));
});

gulp.task('watch', function() {
  gulp.watch('public/stylesheets/*.scss', ['sass']);
  gulp.watch(['public/javascript/*.coffee'], ['coffee', 'compress']);
});

gulp.task('default', ['sass', 'coffee', 'compress', 'fonts', 'css', 'watch']);


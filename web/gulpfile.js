var gulp = require('gulp');
var sass = require('gulp-sass');
var csso = require('gulp-csso');
var uglify = require('gulp-uglify');
var concat = require('gulp-concat');
var plumber = require('gulp-plumber');
var del = require('del');

var paths = {
  fonts: 'bower_components/font-awesome/fonts/*',
  rickshaw: 'bower_components/rickshaw/rickshaw.min.css'
};

gulp.task('clean', function(cb) {
  del(['public/javascript/app.min.js', 'fonts', 'public/stylesheets/*.css'], cb);
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
    'public/javascript/*.js'
  ])
    .pipe(concat('app.min.js'))
    .pipe(gulp.dest('public/javascript'));
});

// Copy all fonts to the public folder
gulp.task('fonts', function() {
  return gulp.src(paths.fonts)
    .pipe(gulp.dest('public/fonts'));
});

// Copy all additional required CSS files to the public folder
gulp.task('css', function() {
  return gulp.src(paths.rickshaw)
    .pipe(gulp.dest('public/stylesheets'));
});

gulp.task('watch', function() {
  gulp.watch('public/stylesheets/*.scss', ['sass']);
  gulp.watch(['public/javascripts/*.js', '!public/app.min.js'], ['compress']);
});

gulp.task('default', ['clean', 'sass', 'compress', 'fonts', 'css', 'watch']);

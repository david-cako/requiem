#!/usr/bin/env node

var fs = require('fs');
var path = require('path');

var ghostJsonPath = process.argv[2];
var outputPath = process.argv[3];

if (!ghostJsonPath || !outputPath) {
    console.log("");        
    console.log("usage:");
    console.log("    ghost-to-jekyll json_dump_path output_path");
    console.log("");
    console.log("Get your blog dump from YOURBLOG.COM/ghost/settings/labs -> 'export'.");
    console.log("");
    console.log("");    
    process.exit(1);
}

var ghostJson = JSON.parse(fs.readFileSync(ghostJsonPath, 'utf8'));

var posts = ghostJson['db'][0]['data']['posts'];

process.chdir(outputPath);

for (i = 0; i < posts.length; i++) {
    var post = posts[i];
    if (post['status'] != 'published') {
        continue;
    }
    var fileName = post['published_at'].split(' ')[0] + '-' + post['slug'] + '.md';
    var stream = fs.createWriteStream(fileName);

    stream.write('---\n');
    stream.write('title: "' + post['title'] + '"\n');
    stream.write('date: "' + post['published_at'] + '"\n');
    stream.write('---\n\n'); 
    stream.write(post['markdown']);
    stream.end();
}
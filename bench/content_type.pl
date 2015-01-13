#!/usr/bin/perl

use strict;
use Data::Dumper;
use Analyze;

if (scalar(@ARGV) < 1) {
    die "Usage:\n\t$0 <performance.csv> [<raw data directory>]\n";
}

open(my $fh, $ARGV[0]);
my $base = $ARGV[1] || '.';
$base =~ s/\///g;

my @lines = <$fh>;
my %performance = map { chomp($_) && split(",", $_) } @lines;

my $content_type_map = Analyze::get_files_by_content_type($base);

Analyze::display_stats($content_type_map, \%performance);

close($fh);

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

my $encoded_words_map = Analyze::get_files_by_encoded_words($base);

Analyze::display_stats($encoded_words_map, \%performance);

close($fh);

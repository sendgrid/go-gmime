#!/usr/bin/perl

use strict;
use Data::Dumper;
use Analyze;

if (scalar(@ARGV) < 1) {
    die "Usage:\n\t$0 <errors.log>\n";
}

open(my $fh, $ARGV[0]);

my @lines = <$fh>;
my @failures = map { chomp($_) && grep { /Failed to parse: (.*)$/; $_ = $1 } $_ } @lines;

my $content_type_map = Analyze::get_files_by_content_type(\@failures);
Analyze::display_stats($content_type_map);
print "\n\n";

my $content_transfer_encoding_map = Analyze::get_files_by_content_transfer_encoding(\@failures);
Analyze::display_stats($content_transfer_encoding_map);
print "\n\n";

my $encoded_words_map = Analyze::get_files_by_encoded_words(\@failures);
Analyze::display_stats($encoded_words_map);
print "\n\n";

close($fh);

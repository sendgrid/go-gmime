#!/usr/bin/perl
# Run
# cd go-gmime/bench
# perl ./perl/test.pl

use strict;
use File::Find::Iterator;
use IO::Handle;
use Time::HiRes;
use MIME::Parser;

my $root = "./data";
my ($count, $time) = parsePerl($root);
print STDERR "Parsed $count files in $time seconds with Perl.\n";

sub parsePerl {
    my $directory = shift;
    my $find = File::Find::Iterator->create(dir => [$directory],
                                            filter => sub { -f && /content-/ });

    my $total_files = 0;
    my $total_time = 0;

    while (my $file = $find->next) {
        my $fh = IO::File->new($file);
        my $parser = new MIME::Parser;
        $parser->output_under("/tmp");
        my $start = Time::HiRes::time();
        my $entity = $parser->parse($fh);
        my $end = Time::HiRes::time();
        my $time = $end - $start;
        if (defined($entity)) {
            print "$file,$time\n";
            $total_files += 1;
            $total_time += $time;
        } else {
            print STDERR "Failed to parse $file\n";
        }
    }

    return $total_files, $total_time;
}
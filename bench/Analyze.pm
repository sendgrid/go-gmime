package Analyze;

use strict;
use Encode;

sub display_stats {
    my $map = shift;
    my $measurements = shift;
    my $infinity = encode("UTF-8", "\x{221E}");

    print "Attribute", ",", "Count", ",", "Total", ",", "Mean", "\n";
    for my $key (keys %$map) {
        my @values = unique(@{$map->{$key}});
        @values = map { $measurements->{$_} } @values if defined($measurements);
        my $count = scalar(@values);
        my $sum = 0;
        foreach my $value (@values) { $sum += $value }
        my $average = $sum / $count;
        $sum = $infinity if $sum == 0;
        $average = $infinity if $average == 0;
        print $key, ",", $count, ",", $sum, ",", $average, "\n";
    }
}

sub get_files_by_content_type {
    my $files = shift;
    return get_files_by_specific_header_value($files, 'Content-Type');
}

sub get_files_by_content_transfer_encoding {
    my $files = shift;
    return get_files_by_specific_header_value($files, 'content-transfer-encoding');
}

sub demangle {
    my $key = lc(shift);
    $key =~ s/[ ;]//g;
    return $key;
}

sub get_files_by_specific_header_value {
    my $files = shift;
    my $specific_header = shift();
    my $search = (ref($files) ne 'ARRAY') ? [$files] : $files;
    $search = join(" ", @$search);
    my @matching_lines = grep { chomp } `grep $search -oREe "^$specific_header:([^;]+);?"`;
    my @matching_parts = map { $_ = [split /^([^:]+):(.+)$/]; shift(@$_); $_ } @matching_lines;
    return make_map_of_tuple_list(\&demangle, \@matching_parts);
}

sub get_files_by_encoded_words {
    my $files = shift;
    my $regex = '=(\?\S+?\?\S\?)\S+\?=';
    my $search = (ref($files) ne 'ARRAY') ? [$files] : $files;
    $search = join(" ", @$search);
    my @matching_lines = grep { chomp } `grep $search -oREe "$regex"`;
    my @matching_parts = map { $_ = [split /^([^:]+):$regex/]; shift(@$_); $_ } @matching_lines;
    return make_map_of_tuple_list(\&demangle, \@matching_parts);
}

sub get_files_by_header {
    my $files = shift;
    my @matching_lines = grep { chomp } `grep $files -oREe "^([A-Z][a-zA-Z\-]+:.*);"`;
    my @matching_parts = map { $_ = [split /^([^:]+):([^:]+):.+/]; shift(@$_); $_ } @matching_lines;
    return make_map_of_tuple_list(\&demangle, \@matching_parts);
}

sub make_map_of_tuple_list {
    my $demangler = shift;
    my $tuple_list = shift;
    my %map;
    for my $tuple (@$tuple_list) {
        my ($value, $mangled_key) = @$tuple;
        my $key = $demangler->($mangled_key);
        next if $key eq '';
        $map{$key} = [] unless defined($map{$key});
        push(@{$map{$key}}, $value);
    }
    return \%map;
}

sub unique {
    my @array = @_;
    my %map = grep { $_ => 1 } @array;
    return keys %map;
}

1;

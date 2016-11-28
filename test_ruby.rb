#!/usr/bin/ruby

contents = IO.binread("challenge.bin")

new_arr = Array.new

contents.bytes.to_a.each_slice(2) do |little, big|
  new_val = big * (2 ** 8) + little
  puts "#{new_val.to_s}"
end

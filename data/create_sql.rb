lines = File.open(ARGV[0]).readlines
desc, tbl, col,ftype = lines[0..3]
tbl = tbl.strip
desc = desc.strip
cols = col.chomp.split("\t")
ftypes = ftype.chomp.split("\t")

File.open("i_#{tbl}.data","w") do |f|
    lines[4..-1].each do |line|
        f.puts line.gsub(",","")
    end
end

File.open("#{tbl}.sql","w") do |f|
    f.puts "DELETE FROM idx where name ='#{tbl}';"
    f.puts "INSERT INTO idx (name, desc) VALUES ('#{tbl}','#{desc}');"
    f.puts "DROP TABLE IF EXISTS #{tbl};" 
    f.puts "CREATE TABLE #{tbl} ("
    (0...cols.length-1).each do |i|
        f.puts "\t#{cols[i]} #{ftypes[i]}, "
    end
    f.puts "\t#{cols[-1]} #{ftypes[-1]}\n);"
    f.puts ".separator \"\\t\""
    f.puts ".import i_#{tbl}.data #{tbl}"
end

`sqlite3 data.db < #{tbl}.sql`

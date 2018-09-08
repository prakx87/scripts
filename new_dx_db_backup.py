#!/usr/bin/python
__author__ = 'prakx'

import datetime
import subprocess
import MySQLdb as mariadb
import dateutil.parser


def run_commands(linux_cmd):
    p = subprocess.Popen(linux_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    out, error = p.communicate()
    returncode = p.returncode
    if returncode != 0:
        raise Exception(error)
    return


def del_cloud_files(cloud_path, del_command):
    cloud_type = cloud_path.split(':')[0]
    dp = subprocess.Popen(del_command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    dout, derror = dp.communicate()
    cleaned_dout = dout.split('\n')
    total_list = []
    for bk in cleaned_dout:
        if bk:
            bk_details = bk.split()
            total_list.append(bk_details[-1].strip())
    for oldbk in total_list:
        dbbk = oldbk.split('-')[0]
        if dbbk in mysql_databases:
            bktime = '-'.join(oldbk.split('-')[1].split('.')[0].split('_')[:-1])
            oldbktime = dateutil.parser.parse(bktime)
            time_between_bks = datetime.datetime.now() - oldbktime
            if time_between_bks.days > 60:
                print "Deleting {0} backup: {1}".format(cloud_type, oldbk)
                del_file_cmd = "rclone delete {0}{1}".format(cloud_path, oldbk)
                run_commands(del_file_cmd)
    return


# MySQL settings and other constants
backup_date = datetime.datetime.now().strftime("%Y_%m_%d_%H")
backup_dir = "/root/backups"
gdrive_path = "gdrive_remote:Doramax264/DB_Backups/"
yandex_path = "yandex:Doramax264/DB_Backups/"

# Create backup directory and set permissions
print "Date:", backup_date, "\n"

# Get MySQL databases
conn = mariadb.connect(host="localhost", read_default_file="/etc/my.cnf.d/backup.cnf")
cur = conn.cursor()
cur.execute("SHOW DATABASES")
l = cur.fetchall()
mysql_databases = [i[0] for i in l]
cur.close()
conn.close()

# Backup and compress each database
bk_taken = []
for database in mysql_databases:
    if database != "information_schema" and database != "performance_schema" and database != "mysql":
        backup_name = "{0}-{1}.gz".format(database, backup_date)
        print "Creating backup of {0} database: {1}".format(database, backup_name)
        bk_taken.append(backup_name)
        cmd1 = "mysqldump --defaults-extra-file=/etc/my.cnf.d/backup.cnf {0} | gzip > {1}/{2}".format(database, backup_dir, backup_name)
        run_commands(cmd1)
        cmd2 = "chmod 600 {0}/{1}".format(backup_dir, backup_name)
        run_commands(cmd2)

# Copy the backup files to Google drive and Yandex
for bk_file in bk_taken:
    print "Copying backup {0} to Cloud drives.".format(bk_file)
    gdrive_cmd = "rclone copy /root/backups/{0} {1} -v --log-file /root/scripts/backup_log.txt".format(bk_file,
                                                                                                       gdrive_path)
    yandex_cmd = "rclone copy {0}{1} {2} -v --log-file /root/scripts/backup_log.txt".format(gdrive_path, bk_file,
                                                                                            yandex_path)
    run_commands(gdrive_cmd)
    run_commands(yandex_cmd)

# Delete backup files older than 60 days in Cloud
gdrive_del = "rclone lsl {0}".format(gdrive_path)
yandex_del = "rclone lsl {0}".format(yandex_path)
del_cloud_files(gdrive_path, gdrive_del)
del_cloud_files(yandex_path, yandex_del)

# Delete backup files older than 10 days in VPS
del_cmd = "find /root/backups -mtime +10 -exec rm {} \;"
run_commands(del_cmd)

#print "\n======================================================================================"
